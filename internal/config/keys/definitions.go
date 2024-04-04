package keys

var DryRun = newKey("dry-run", "dry-run", true)
var Debug = newKey("debug", "debug", true)

// backup/restore shared
var BackupId = newKey("backup-id", "backup id to use", "")
var BackupDir = newKey("backup-dir", "directory to store backups", "")
var DockerHosts = newKey("docker-hosts", "docker hosts to perform actions on", []string{"unix:///var/run/docker.sock"})
var LabelPrefix = newKey("label-prefix", "label prefix used to scan containers", "bckupr")

var NoDaemon = newKey("no-daemon", "don't use bckupr daemon, metrics won't be available using this flag", true)
var DaemonNet = newKey("daemon-net", "don't use bckupr daemon, metrics won't be available using this flag", "unix")
var DaemonProtocol = newKey("daemon-protocol", "don't use bckupr daemon, metrics won't be available using this flag", "http")
var DaemonAddr = newKey("daemon-addr", "don't use bckupr daemon, metrics won't be available using this flag", UnixSocket.Default)

var LocalContainers = newKey("local-containers", "container config for managing local backups", "./configs/local/tar.yml")
var OffsiteContainers = newKey("offsite-containers", "container config for managing offsite backups", "")

var NotificationUrls = newKey("notification-urls", "shoutrrr notification urls", []string{
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
var IncludeNames = newKey("include-names", "names of containers to include", []string{})
var IncludeVolumes = newKey("include-volumes", "names of volumes to include", []string{})
var ExcludeName = newKey("exclude-names", "names of containers to exclude", []string{})
var ExcludeVolumes = newKey("exclude-volumes", "names of volumes to exclude", []string{})

// cron
var TimeZone = newKey("timezone", "cron timezeon", "UTC")
var BackupSchedule = newKey("backup-schedule", "cron schedule for backups", "0 0 * * *")

// web
var UnixSocket = newKey("unix-socket", "unix socket to bind web server", ".bckupr.sock")
var TcpAddr = newKey("tcp-addr", "tcp address to bind web server", "0.0.0.0:8000")
var ExposeApi = newKey("expose-api", "exposes /api endpoints via tcp", false)
var UiEnabled = newKey("ui-enabled", "exposes ui", true)
var MetricsEnabled = newKey("metrics-enabled", "enables /metrics endpoint", false)

// END OF DEFINITIONS
func stopModes(stopModes []string) *Key {
	return newKey("stop-modes", "stop modes used to controller shutdown: all, labelled, linked, writers, attached", stopModes)
}
