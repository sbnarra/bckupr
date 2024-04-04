package encodings

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v3"
)

func FromYaml(reader io.Reader, data any) error {
	if err := yaml.NewDecoder(reader).Decode(data); err != nil {
		return err
	}
	return nil
}

func ToYaml(data any) (string, error) {
	buffer := bytes.NewBuffer([]byte{})
	if err := yaml.NewEncoder(buffer).Encode(data); err != nil {
		return "", err
	} else {
		return buffer.String(), nil
	}
}
