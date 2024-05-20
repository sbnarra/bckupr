package server

import (
	"context"
	"net/http"

	"github.com/sbnarra/bckupr/internal/notifications"
)

type server struct {
	*http.Server
	context.Context
	Config
}

type Config struct {
	ContainerBackupDir string
	HostBackupDir      string
	DockerHosts        []string

	TcpAddr string

	LocalContainersConfig   string
	OffsiteContainersConfig string

	NotificationSettings *notifications.NotificationSettings
}
