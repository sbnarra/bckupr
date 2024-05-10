package notifications

type NotificationSettings struct {
	NotificationUrls    []string `json:"notification-urls"`
	NotifyJobStarted    bool     `json:"notify-job-started"`
	NotifyJobCompleted  bool     `json:"notify-job-completed"`
	NotifyJobError      bool     `json:"notify-job-error"`
	NotifyTaskStarted   bool     `json:"notify-task-started"`
	NotifyTaskCompleted bool     `json:"notify-task-completed"`
	NotifyTaskError     bool     `json:"notify-task-error"`
}