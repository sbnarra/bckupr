package types

import "time"

type CreateBackupRequest struct {
	Args                 TaskArgs              `json:"args"`
	NotificationSettings *NotificationSettings `json:"notification-settings"`
}

type RotateBackupsRequest struct {
	Destroy      bool   `json:"destroy"`
	PoliciesPath string `json:"policies-path"`
}

type RestoreBackupRequest struct {
	Args                 TaskArgs              `json:"args"`
	NotificationSettings *NotificationSettings `json:"notification-settings"`
}

type DaemonInput struct {
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

type TaskArgs struct {
	BackupId    string  `json:"backup-id"`
	LabelPrefix string  `json:"label-prefix"`
	Filters     Filters `json:"filters"`
}

type NotificationSettings struct {
	NotificationUrls    []string `json:"notification-urls"`
	NotifyJobStarted    bool     `json:"notify-job-started"`
	NotifyJobCompleted  bool     `json:"notify-job-completed"`
	NotifyJobError      bool     `json:"notify-job-error"`
	NotifyTaskStarted   bool     `json:"notify-task-started"`
	NotifyTaskCompleted bool     `json:"notify-task-completed"`
	NotifyTaskError     bool     `json:"notify-task-error"`
}

type Filters struct {
	StopModes      []string `json:"stop-modes"`
	IncludeNames   []string `json:"include-names"`
	IncludeVolumes []string `json:"include-volumes"`
	ExcludeNames   []string `json:"exclude-names"`
	ExcludeVolumes []string `json:"exclude-volumes"`
}

type Volume struct {
	Name    string    `json:"name"`
	Ext     string    `json:"ext"`
	Mount   string    `json:"mount"`
	Created time.Time `json:"created"`
	Size    int64     `json:"size(kb)"`
	Error   string    `json:"error"`
}

type Backup struct {
	Id      string    `json:"backup-id"`
	Type    string    `json:"type"`
	Created time.Time `json:"created"`
	Volumes []Volume  `json:"volumes"`
}

type ContainerTemplates struct {
	Local   LocalContainerTemplates
	Offsite *OffsiteContainerTemplates
}

type LocalContainerTemplates struct {
	Backup  ContainerTemplate `yaml:"backup"`
	Restore ContainerTemplate `yaml:"restore"`
	FileExt string            `yaml:"file-ext"`
}

type OffsiteContainerTemplates struct {
	OffsitePush ContainerTemplate `yaml:"offsite-push"`
	OffsitePull ContainerTemplate `yaml:"offsite-pull"`
}

type ContainerTemplate struct {
	Image   string            `yaml:"image"`
	Cmd     []string          `yaml:"cmd"`
	Env     []string          `yaml:"env"`
	Volumes []string          `yaml:"volumes"`
	Labels  map[string]string `yaml:"labels"`
}
