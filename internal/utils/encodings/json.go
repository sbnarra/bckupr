package encodings

import (
	"encoding/json"

	"github.com/sbnarra/bckupr/internal/utils/errors"
)

func ToJsonIE(data any) string { // ignore error
	j, _ := ToJson(data)
	return j
}

func ToJson(data any) (string, *errors.Error) {
	if b, err := json.MarshalIndent(data, "", "  "); err != nil {
		return "", errors.Wrap(err, "error encoding to json")
	} else {
		return string(b), nil
	}
}
