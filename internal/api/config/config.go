package config

import "github.com/sbnarra/bckupr/internal/config/keys"

type Config struct {
	HostBackupDir string
	DockerHosts   []string

	// UnixSocket string
	TcpAddr string

	LocalContainersConfig   string
	OffsiteContainersConfig string
}

func New() Config {
	return Config{
		DockerHosts:   keys.DockerHosts.EnvStringSlice(),
		HostBackupDir: keys.HostBackupDir.EnvString(),

		LocalContainersConfig:   keys.LocalContainersConfig.EnvString(),
		OffsiteContainersConfig: keys.OffsiteContainersConfig.EnvString(),

		// UnixSocket: keys.UnixSocket.EnvString(),
		TcpAddr: keys.TcpAddr.EnvString(),
	}
}
