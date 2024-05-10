package server

import "github.com/sbnarra/bckupr/internal/notifications"

type Config struct {
	HostBackupDir string
	DockerHosts   []string

	// UnixSocket string
	TcpAddr string

	LocalContainersConfig   string
	OffsiteContainersConfig string

	NotificationSettings *notifications.NotificationSettings
}
