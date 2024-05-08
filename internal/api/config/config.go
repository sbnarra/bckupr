package config

import "github.com/sbnarra/bckupr/internal/config/keys"

type Config struct {
	BackupDir               string   `json:"backup-dir"`
	DockerHosts             []string `json:"docker-hosts"`
	LocalContainersConfig   string   `json:"local-containers-config"`
	OffsiteContainersConfig string   `json:"offsite-containers-config"`

	UnixSocket string `json:"unix-socket"`
	TcpAddr    string `json:"tcp-addr"`
	TcpApi     bool   `json:"tcp-api"`
	UI         bool   `json:"ui-enabled"`
	Metrics    bool   `json:"metrics-enabled"`
}

func New() Config {
	return Config{
		BackupDir:               keys.HostBackupDir.EnvString(),
		LocalContainersConfig:   keys.LocalContainersConfig.EnvString(),
		OffsiteContainersConfig: keys.OffsiteContainersConfig.EnvString(),

		UnixSocket: keys.UnixSocket.EnvString(),
		TcpAddr:    keys.TcpAddr.EnvString(),
		TcpApi:     keys.TcpApi.EnvBool(),
		UI:         keys.UI.EnvBool(),
		Metrics:    keys.Metrics.EnvBool(),
	}
}
