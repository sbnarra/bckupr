package server

import "github.com/sbnarra/bckupr/internal/notifications"

type Config struct {
	ContainerBackupDir string
	HostBackupDir      string
	DockerHosts        []string

	TcpAddr string

	LocalContainersConfig   string
	OffsiteContainersConfig string

	NotificationSettings *notifications.NotificationSettings
}
