package meta

import (
	"github.com/sbnarra/bckupr/internal/api/spec"
)

type storage struct {
	data map[string]*spec.Backup
}
