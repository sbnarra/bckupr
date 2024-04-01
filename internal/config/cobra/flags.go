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
		flagSet.String(key.CliId, key.Default.(string), key.Desc)
	} else if typ == "int" {
		flagSet.Int(key.CliId, key.Default.(int), key.Desc)
	} else if typ == "bool" {
		flagSet.Bool(key.CliId, key.Default.(bool), key.Desc)
	} else if typ == "[]string" {
		flagSet.StringSlice(key.CliId, key.Default.([]string), key.Desc)
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
	if err == nil {
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
	if err == nil {
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
