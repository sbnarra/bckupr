package rotate

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/meta/reader"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/xhit/go-str2duration/v2"
)

type Period struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type RotatePolicy struct {
	Period Period `json:"period"`
	Keep   int    `json:"keep"`
}

type RotatePolicies struct {
	Policies []RotatePolicy `json:"policies"`
}

func Rotate(ctx contexts.Context, input spec.RotateInput) *errors.Error {
	policies := &RotatePolicies{}
	if _, err := os.Stat(input.PoliciesPath); err == nil {
		if handle, err := os.Open(input.PoliciesPath); err != nil {
			return errors.Wrap(err, "failed to read: "+input.PoliciesPath)
		} else if err := encodings.FromYaml(bufio.NewReader(handle), policies); err != nil {
			return err
		}
	}

	if reader, err := reader.Load(ctx); err != nil {
		return err
	} else if err := sortAndValidate(policies.Policies); err != nil {
		return err
	} else {
		for _, policy := range policies.Policies {
			if err := applyPolicy(ctx, input.Destroy, policy, reader); err != nil {
				logging.CheckError(ctx, err, "error applying rotation policy")
			}
		}
		return nil
	}
}

func sortAndValidate(policies []RotatePolicy) *errors.Error {
	sort.Slice(policies, func(i, j int) bool {
		return policies[i].Period.From < policies[j].Period.From
	})

	var lastPolicy *RotatePolicy
	for _, policy := range policies {
		if policy.Keep == 0 {
			return errors.New("invalid policy, keeps 0 backups")
		}

		if lastPolicy != nil {
			if _, lastPolicyEnd, err := parsePeriod(lastPolicy.Period); err != nil {
				return err
			} else if policyStart, _, err := parsePeriod(policy.Period); err != nil {
				return err
			} else if policyStart.After(lastPolicyEnd) {
				return errors.Errorf("invalid policy periods: last policy to: %v < policy from: %v", lastPolicy.Period.To, policy.Period.From)
			}
		}
		lastPolicy = &policy
	}
	return nil
}

func parsePeriod(period Period) (time.Time, time.Time, *errors.Error) {
	policyStart := time.Now()
	policyEnd := time.Now()
	if fromDuration, err := str2duration.ParseDuration(period.From); err != nil {
		return policyStart, policyEnd, errors.Wrap(err, "failed to parse: "+period.From)
	} else if toDuration, err := str2duration.ParseDuration(period.To); err != nil {
		return policyStart, policyEnd, errors.Wrap(err, "failed to parse: "+period.To)
	} else {
		policyStart = policyStart.Add(fromDuration)
		policyEnd = policyEnd.Add(toDuration)
		return policyStart, policyEnd, nil
	}
}

func applyPolicy(ctx contexts.Context, destroyBackups bool, policy RotatePolicy, reader *reader.Reader) *errors.Error {
	if policyStart, policyEnd, err := parsePeriod(policy.Period); err != nil {
		return err
	} else {
		backups := []*spec.Backup{}
		for _, b := range reader.Find() {
			if b.Created.After(policyStart) && b.Created.Before(policyEnd) {
				backups = append(backups, b)
			}
		}

		if len(backups) == 0 {
			logging.Info(ctx, fmt.Sprintf("no backups in period: (%v) %v  <->  (%v) %v",
				policy.Period.From, policyStart.Format("2006-01-02 15:04:05"),
				policy.Period.To, policyEnd.Format("2006-01-02 15:04:05")))
			return nil
		}

		sort.Slice(backups, func(i, j int) bool {
			return backups[i].Created.Before(backups[j].Created)
		})

		if policy.Keep > 0 {
			if len(backups) > policy.Keep {
				newest := len(backups) - policy.Keep
				rotateBackups(ctx, destroyBackups, backups[:newest])
			} else {
				logging.Info(ctx, "no backups to rotate")
			}
		} else {
			if len(backups) > (policy.Keep * -1) {
				oldest := policy.Keep * -1
				rotateBackups(ctx, destroyBackups, backups[oldest:])
			} else {
				logging.Info(ctx, "no backups to rotate")
			}
		}
		return nil
	}
}

func rotateBackups(ctx contexts.Context, destroyBackups bool, backups []*spec.Backup) {
	binPath := filepath.Join(ctx.HostBackupDir, "bin")
	for _, backup := range backups {
		backupPath := filepath.Join(ctx.ContainerBackupDir, backup.Id)
		if err := rotateBackup(ctx, destroyBackups, backupPath, filepath.Join(binPath, backup.Id)); err != nil {
			logging.CheckError(ctx, err)
		}
	}
}

func rotateBackup(ctx contexts.Context, destroyBackups bool, backupPath string, binPath string) *errors.Error {
	if !destroyBackups {
		logging.Info(ctx, "ln", backupPath, binPath)
		if !ctx.DryRun {
			if err := os.Link(backupPath, binPath); err != nil {
				return errors.Wrap(err, "failed to link "+backupPath+" to "+binPath)
			}
		}
	}

	logging.Info(ctx, "rm -rf", backupPath)
	if !ctx.DryRun {
		err := os.RemoveAll(backupPath)
		return errors.Wrap(err, "failed to remove "+backupPath)
	}
	return nil
}
