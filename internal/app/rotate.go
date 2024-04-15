package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sbnarra/bckupr/internal/meta"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/logging"
	"github.com/sbnarra/bckupr/pkg/types"
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

func Rotate(ctx contexts.Context, input *types.RotateBackupsRequest) error {
	policies := &RotatePolicies{}
	if _, err := os.Stat(input.PoliciesPath); err == nil {
		if handle, err := os.Open(input.PoliciesPath); err != nil {
			return err
		} else if err := encodings.FromYaml(bufio.NewReader(handle), policies); err != nil {
			return err
		}
	}

	if reader, err := meta.NewReader(ctx); err != nil {
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

func sortAndValidate(policies []RotatePolicy) error {
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
				return fmt.Errorf("invalid policy periods: last policy to: %v < policy from: %v", lastPolicy.Period.To, policy.Period.From)
			}
		}
		lastPolicy = &policy
	}
	return nil
}

func parsePeriod(period Period) (time.Time, time.Time, error) {
	policyStart := time.Now()
	policyEnd := time.Now()
	if fromDuration, err := str2duration.ParseDuration(period.From); err != nil {
		return policyStart, policyEnd, err
	} else if toDuration, err := str2duration.ParseDuration(period.To); err != nil {
		return policyStart, policyEnd, err
	} else {
		policyStart = policyStart.Add(fromDuration)
		policyEnd = policyEnd.Add(toDuration)
		return policyStart, policyEnd, nil
	}
}

func applyPolicy(ctx contexts.Context, destroyBackups bool, policy RotatePolicy, reader meta.Reader) error {
	if policyStart, policyEnd, err := parsePeriod(policy.Period); err != nil {
		return err
	} else {
		backups := []*types.Backup{}
		reader.ForEach(func(b *types.Backup) error {
			if b.Created.After(policyStart) && b.Created.Before(policyEnd) {
				backups = append(backups, b)
			}
			return nil
		})

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

func rotateBackups(ctx contexts.Context, destroyBackups bool, backups []*types.Backup) {
	binPath := filepath.Join(ctx.BackupDir, "bin")
	for _, backup := range backups {
		backupPath := filepath.Join(ctx.BackupDir, backup.Id)
		if err := rotateBackup(ctx, destroyBackups, backupPath, filepath.Join(binPath, backup.Id)); err != nil {
			logging.CheckError(ctx, err)
		}
	}
}

func rotateBackup(ctx contexts.Context, destroyBackups bool, backupPath string, binPath string) error {
	if !destroyBackups {
		logging.Info(ctx, "ln", backupPath, binPath)
		if !ctx.DryRun {
			if err := os.Link(backupPath, binPath); err != nil {
				return err
			}
		}
	}

	logging.Info(ctx, "rm -rf", backupPath)
	if !ctx.DryRun {
		return os.RemoveAll(backupPath)
	}
	return nil
}