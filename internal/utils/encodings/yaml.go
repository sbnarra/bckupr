package encodings

import (
	"io"

	"gopkg.in/yaml.v3"
)

func FromYaml(reader io.Reader, data any) error {
	if err := yaml.NewDecoder(reader).Decode(data); err != nil {
		return err
	}
	return nil
}
