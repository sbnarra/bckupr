package containers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/sbnarra/bckupr/internal/utils/encodings"
	"github.com/sbnarra/bckupr/pkg/types"
)

func ContainerTemplates(local string, offsite string) (types.LocalContainerTemplates, *types.OffsiteContainerTemplates, error) {
	if local, err := LocalContainerTemplates(local); err != nil {
		return local, nil, err
	} else if offsite, err := OffsiteContainerTemplates(offsite); err != nil {
		return local, nil, err
	} else {
		return local, offsite, nil
	}
}

func LocalContainerTemplates(location string) (types.LocalContainerTemplates, error) {
	config := types.LocalContainerTemplates{}
	err := loadContainerTemplates(location, "local", &config)
	return config, err
}

func OffsiteContainerTemplates(location string) (*types.OffsiteContainerTemplates, error) {
	if location == "" {
		return nil, nil
	}
	config := &types.OffsiteContainerTemplates{}
	err := loadContainerTemplates(location, "offsite", config)
	return config, err
}

func loadContainerTemplates[T any](location string, usage string, data T) error {
	if reader, err := getReader(location, usage); err != nil {
		return err
	} else {
		return encodings.FromYaml(reader, data)
	}
}

func getReader(location string, usage string) (io.Reader, error) {
	if parsed, err := url.Parse(location); err != nil {
		return nil, err
	} else {
		switch parsed.Scheme {
		case "file", "":
			return fileRead(parsed.Path)
		case "http", "https":
			return httpGet(location)
		default:
			return nil, fmt.Errorf("unsupported scheme for %v containers: '%v'", usage, parsed.Scheme)
		}
	}
}

func httpGet(location string) (io.Reader, error) {
	if res, err := http.Get(location); err != nil {
		return nil, err
	} else {
		return res.Body, nil
	}
}

func fileRead(location string) (io.Reader, error) {
	if reader, err := os.Open(location); err != nil {
		return nil, err
	} else {
		return reader, err
	}
}
