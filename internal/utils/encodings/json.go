package encodings

import (
	"encoding/json"
)

func ToJson_(data any) string {
	j, _ := ToJson(data)
	return j
}

func ToJson(data any) (string, error) {
	if b, err := json.MarshalIndent(data, "", "  "); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
