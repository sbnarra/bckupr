package cobra

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sbnarra/bckupr/internal/config/keys"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func register(key *keys.Key, flagSet *pflag.FlagSet) {
	typ := fmt.Sprintf("%T", key.Default)
	if typ == "string" {
		defaultV := key.Default.(string)
		if key.EnvExists() {
			defaultV = key.EnvString()
		}
		flagSet.String(key.CliId, defaultV, key.Desc)
	} else if typ == "int" {
		defaultV := key.Default.(int)
		if key.EnvExists() {
			defaultV = key.EnvInt()
		}
		flagSet.Int(key.CliId, defaultV, key.Desc)
	} else if typ == "bool" {
		defaultV := key.Default.(bool)
		if key.EnvExists() {
			defaultV = key.EnvBool()
		}
		flagSet.Bool(key.CliId, defaultV, key.Desc)
	} else if typ == "[]string" {
		defaultV := key.Default.([]string)
		if key.EnvExists() {
			defaultV = key.EnvStringSlice()
		}
		flagSet.StringSlice(key.CliId, defaultV, key.Desc)
	} else {
		panic("unknown key type: " + typ + ": " + key.Id)
	}
}

func required(key *keys.Key, cmd *cobra.Command) {
	if !key.EnvExists() {
		cmd.MarkFlagRequired(key.CliId)
	}
}

func String(key *keys.Key, flags *pflag.FlagSet) (string, error) {
	val, err := flags.GetString(key.CliId)
	if err == nil && len(val) != 0 {
		os.Setenv(key.EnvId(), val)
	}
	return val, err
}

func Int(key *keys.Key, flags *pflag.FlagSet) (int, error) {
	val, err := flags.GetInt(key.CliId)
	if err == nil {
		os.Setenv(key.EnvId(), strconv.Itoa(val))
	}
	return val, err
}

func StringSlice(key *keys.Key, flags *pflag.FlagSet) ([]string, error) {
	val, err := flags.GetStringSlice(key.CliId)
	if err == nil && len(val) != 0 {
		os.Setenv(key.EnvId(), strings.Join(val, ","))
	}
	return val, err
}

func Bool(key *keys.Key, flags *pflag.FlagSet) (bool, error) {
	val, err := flags.GetBool(key.CliId)
	if err == nil {
		os.Setenv(key.EnvId(), strconv.FormatBool(val))
	}
	return val, err
}
