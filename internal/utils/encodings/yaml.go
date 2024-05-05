package encodings

import (
	"bytes"
	"io"

	"github.com/sbnarra/bckupr/internal/utils/errors"
	"gopkg.in/yaml.v3"
)

func FromYaml(reader io.Reader, data any) *errors.Error {
	if err := yaml.NewDecoder(reader).Decode(data); err != nil {
		return errors.Wrap(err, "error decoding from yaml")
	}
	return nil
}

func ToYaml(data any) (string, *errors.Error) {
	buffer := bytes.NewBuffer([]byte{})
	if err := yaml.NewEncoder(buffer).Encode(data); err != nil {
		return "", errors.Wrap(err, "error encoding to yaml")
	} else {
		return buffer.String(), nil
	}
}
