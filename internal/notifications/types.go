package notifications

import "github.com/sbnarra/bckupr/internal/config/keys"

type NotificationSettings struct {
	NotificationUrls    []string `json:"notification-urls"`
	NotifyJobStarted    bool     `json:"notify-job-started"`
	NotifyJobCompleted  bool     `json:"notify-job-completed"`
	NotifyJobError      bool     `json:"notify-job-error"`
	NotifyTaskStarted   bool     `json:"notify-task-started"`
	NotifyTaskCompleted bool     `json:"notify-task-completed"`
	NotifyTaskError     bool     `json:"notify-task-error"`
}

func settings() *NotificationSettings {
	return &NotificationSettings{
		NotificationUrls:    keys.NotificationUrls.EnvStringSlice(),
		NotifyJobStarted:    keys.NotifyJobStarted.EnvBool(),
		NotifyJobCompleted:  keys.NotifyJobCompleted.EnvBool(),
		NotifyJobError:      keys.NotifyJobError.EnvBool(),
		NotifyTaskStarted:   keys.NotifyTaskStarted.EnvBool(),
		NotifyTaskCompleted: keys.NotifyTaskCompleted.EnvBool(),
		NotifyTaskError:     keys.NotifyTaskError.EnvBool(),
	}
}
