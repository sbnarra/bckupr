package version

import (
	"os"

	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func Version(ctx contexts.Context) spec.Version {
	return spec.Version{
		Created: os.Getenv("CREATED"),
		Version: os.Getenv("VERSION"),
	}
}
