package keys

import (
	"os"
	"strconv"
	"strings"
)

type Key struct {
	Id      string
	CliId   string
	Desc    string
	Default any
}

func newKey(id string, desc string, defaultV any) *Key {
	return &Key{
		Id:      id,
		CliId:   strings.ReplaceAll(id, ".", "-"),
		Desc:    desc,
		Default: defaultV,
	}
}

func (key *Key) EnvExists() bool {
	_, found := os.LookupEnv(key.EnvId())
	return found
}

func (key *Key) EnvBool() bool {
	if str, found := os.LookupEnv(key.EnvId()); !found {
		return key.Default.(bool)
	} else if b, err := strconv.ParseBool(str); err != nil {
		return key.Default.(bool)
	} else {
		return b
	}
}

func (key *Key) EnvString() string {
	if str, found := os.LookupEnv(key.EnvId()); !found {
		return key.Default.(string)
	} else {
		return str
	}
}

func (key *Key) EnvStringSlice() []string {
	if str, found := os.LookupEnv(key.EnvId()); !found {
		return key.Default.([]string)
	} else {
		res := []string{}
		for _, v := range strings.Split(str, ",") {
			if v != "" {
				res = append(res, v)
			}
		}
		return res
	}
}

func (key *Key) EnvId() string {
	envKey := strings.ReplaceAll(key.Id, ".", "-")
	envKey = strings.ReplaceAll(envKey, "-", "_")
	return strings.ToUpper(envKey)
}
