package keys

var DryRun = newKey("dry-run", "only log actions and not executed", true)
var Debug = newKey("debug", "enable additional logging information", false)

// backup/restore shared
var BackupId = newKey("backup-id", "backup id to use (backup tasks will autogenerate if not set)", "")
var BackupDir = newKey("backup-dir", "backups archieve directory", "")
var DockerHosts = newKey("docker-hosts", "docker host uri's to manage backups on", []string{"unix:///var/run/docker.sock"})
var LabelPrefix = newKey("label-prefix", "label prefix used to scan containers for configuration", "bckupr")

var DaemonNet = newKey("daemon-net", "network connection type for bckupr daemon, unix or tcp", "unix")
var DaemonProtocol = newKey("daemon-protocol", "protocol for bckupr daemon (don't recommend changing)", "http")
var DaemonAddr = newKey("daemon-addr", "bind address for bckupr daemon, should use unix:///tmp/.bckupr or tcp binding like 0.0.0.0:8000", UnixSocket.Default)

var LocalContainers = newKey("local-containers-config", "yaml config for managing local backups", "./configs/local/tar.yml")
var OffsiteContainers = newKey("offsite-containers-config", "yaml config for managing offsite backups", "")

var NotificationUrls = newKey("notification-urls", "shoutrrr service notification urls (https://containrrr.dev/shoutrrr/latest/services/overview/)", []string{
	// "discord://IT4QcVejlF8P5On9Fn6XTJCpjwnkEWhPnV97JI_KJ3ztKuk7aSLc40jK9bu3OeaSowV9@1221065822185853078",
})
var NotifyJobStarted = newKey("notify-job-started", "enable notifications when backups or restores start", false)
var NotifyJobCompleted = newKey("notify-job-completed", "enable notifications when backups or restores complete", false)
var NotifyJobError = newKey("notify-job-error", "enable notifications when backups or restores complete with an error", false)
var NotifyTaskStarted = newKey("notify-task-started", "enable notifications when backups or restores start on specific volumes", false)
var NotifyTaskCompleted = newKey("notify-task-completed", "enable notifications when backups or restores complete on specific volumes", false)
var NotifyTaskError = newKey("notify-task-errors", "enable notifications when backups or restores complete with an error on specific volumes", false)

// ...filters
var BackupStopModes = stopModes([]string{"labelled", "linked", "writers"})
var RestoreStopModes = stopModes([]string{"labelled", "linked", "attached"})
var IncludeNames = newKey("include-names", "only include containers with matching names", []string{})
var IncludeVolumes = newKey("include-volumes", "only include containers with matching volumes", []string{})
var ExcludeName = newKey("exclude-names", "exclude containers with matching names", []string{})
var ExcludeVolumes = newKey("exclude-volumes", "exclude containers with matching volumes", []string{})

// cron
var TimeZone = newKey("timezone", "timezone to use for cron scheduling", "UTC")
var BackupSchedule = newKey("backup-schedule", "cron expression for backups schedule", "0 0 * * *")
var RotateSchedule = newKey("rotate-schedule", "cron expression for rotations schedule", "")

// web
var UnixSocket = newKey("unix-socket", "unix socket to bind daemon", "/tmp/.bckupr.sock")
var TcpAddr = newKey("tcp-addr", "tcp address to bind ui/api", "0.0.0.0:8000")
var ExposeApi = newKey("expose-api", "exposes api via tcp (by default only unix for local connections)", false)
var UiEnabled = newKey("ui-enabled", "exposes gui", true)
var MetricsEnabled = newKey("metrics-enabled", "enables /metrics endpoint", false)

// rotate
var DestroyBackups = newKey("destroy-backups", "destroy backups instead of moving to bin directory", false)
var PoliciesPath = newKey("rotation-policies-config", "rotation policies yaml path", "./configs/rotation/policies.yaml")

// END OF DEFINITIONS
func stopModes(stopModes []string) *Key {
	return newKey("stop-modes", "stop modes to control shutdown targets: all, labelled, linked, writers, attached", stopModes)
}
