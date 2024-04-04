package types

import "time"

type CreateBackupRequest struct {
	DryRun               bool                  `json:"dry-run"`
	Args                 TaskArgs              `json:"args"`
	NotificationSettings *NotificationSettings `json:"notification-settings"`
}

type DeleteBackupRequest struct {
	Args TaskArgs `json:"args"`
}

type RestoreBackupRequest struct {
	DryRun               bool                  `json:"dry-run"`
	Args                 TaskArgs              `json:"args"`
	NotificationSettings *NotificationSettings `json:"notification-settings"`
}

type TaskArgs struct {
	BackupId                string   `json:"backup-id"`
	DockerHosts             []string `json:"docker-hosts"`
	LabelPrefix             string   `json:"label-prefix"`
	Filters                 Filters  `json:"filters"`
	LocalContainersConfig   string   `json:"local-containers-config"`
	OffsiteContainersConfig string   `json:"offsite-containers-config"`
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

type LocalContainerTemplates struct {
	FileExt string            `json:"file-ext"`
	Backup  ContainerTemplate `json:"backup"`
	Restore ContainerTemplate `json:"restore"`
}

type OffsiteContainerTemplates struct {
	OffsitePush ContainerTemplate `json:"offsite-push"`
	OffsitePull ContainerTemplate `json:"offsite-pull"`
}

type ContainerTemplate struct {
	Alias   string   `json:"alias"`
	Image   string   `json:"image"`
	Cmd     []string `json:"cmd"`
	Env     []string `json:"env"`
	Volumes []string `json:"volumes"`
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

type WebInput struct {
	UnixSocket     string
	TcpAddr        string
	ExposeApi      bool
	UiEnabled      bool
	MetricsEnabled bool
}
