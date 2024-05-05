package containers

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/pkg/types"
)

func ContainerTemplates(local string, offsite string) (types.ContainerTemplates, *errors.Error) {
	if local, err := LocalContainerTemplates(local); err != nil {
		return types.ContainerTemplates{}, err
	} else if offsite, err := OffsiteContainerTemplates(offsite); err != nil {
		return types.ContainerTemplates{}, err
	} else {
		return types.ContainerTemplates{
			Local:   local,
			Offsite: offsite,
		}, nil
	}
}

func LocalContainerTemplates(location string) (types.LocalContainerTemplates, *errors.Error) {
	config := types.LocalContainerTemplates{}
	err := loadContainerTemplates(location, "local", &config)
	return config, err
}

func OffsiteContainerTemplates(location string) (*types.OffsiteContainerTemplates, *errors.Error) {
	if location == "" {
		return nil, nil
	}
	config := &types.OffsiteContainerTemplates{}
	err := loadContainerTemplates(location, "offsite", config)
	return config, err
}

func loadContainerTemplates[T any](location string, usage string, data T) *errors.Error {
	if reader, err := getReader(location, usage); err != nil {
		return err
	} else {
		return encodings.FromYaml(reader, data)
	}
}

func getReader(location string, usage string) (io.Reader, *errors.Error) {
	if parsed, err := url.Parse(location); err != nil {
		return nil, errors.Wrap(err, "failed to parse: "+location)
	} else {
		switch parsed.Scheme {
		case "file", "":
			return fileRead(parsed.Path)
		case "http", "https":
			return httpGet(location)
		default:
			return nil, errors.Errorf("unsupported scheme for %v containers: '%v'", usage, parsed.Scheme)
		}
	}
}

func httpGet(location string) (io.Reader, *errors.Error) {
	if res, err := http.Get(location); err != nil {
		return nil, errors.Wrap(err, "failed to GET: "+location)
	} else {
		return res.Body, nil
	}
}

func fileRead(location string) (io.Reader, *errors.Error) {
	if reader, err := os.Open(location); err != nil {
		return nil, errors.Wrap(err, "failed to read file: "+location)
	} else {
		return reader, nil
	}
}
