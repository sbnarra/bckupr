package types

import (
	"sync"

	"github.com/sbnarra/bckupr/internal/docker/client"
)

type Containers struct {
	client      client.DockerClient
	labelPrefix string
}

type Container struct {
	Lock sync.Mutex
	Id   string
	Name string

	Compose Compose

	Running    bool
	WasRunning bool
	Volumes    map[string]ContainerVolume
	Backup     BackupConfig

	Dependancies Dependancies
	Linked/*Children*/ []*Container
}

type Compose struct {
	Project string
	Service string
}

type ContainerVolume struct {
	Writer bool
}

type BackupConfig struct {
	Ignore     bool
	Stop       bool
	Filesystem bool
	Volumes    map[string]string
}

type Dependancies struct {
	Services   []string
	Containers []string
}
