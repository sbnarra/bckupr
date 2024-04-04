package encodings

import (
	"encoding/json"
)

func ToJson(data any) (string, error) {
	if b, err := json.MarshalIndent(data, "", "  "); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
