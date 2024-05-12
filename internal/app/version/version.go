package version

import (
	"context"
	"os"

	"github.com/sbnarra/bckupr/internal/api/spec"
)

func Version(ctx context.Context) spec.Version {
	return spec.Version{
		Created: os.Getenv("CREATED"),
		Version: os.Getenv("VERSION"),
	}
}
