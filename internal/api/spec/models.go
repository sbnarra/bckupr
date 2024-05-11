// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package spec

import (
	"time"
)

// Defines values for Status.
const (
	StatusCompleted Status = "completed"
	StatusError     Status = "error"
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
)

// Defines values for StopModes.
const (
	All      StopModes = "all"
	Attached StopModes = "attached"
	Labelled StopModes = "labelled"
	Linked   StopModes = "linked"
	Writers  StopModes = "writers"
)

// Backup defines model for Backup.
type Backup struct {
	Created time.Time `json:"created"`
	Error   *string   `json:"error,omitempty"`
	Id      string    `json:"id"`
	Status  Status    `json:"status"`
	Type    string    `json:"type"`
	Volumes []Volume  `json:"volumes"`
}

// Error defines model for Error.
type Error struct {
	Error string `json:"error"`
}

// Filters defines model for Filters.
type Filters struct {
	ExcludeNames   []string `json:"exclude_names"`
	ExcludeVolumes []string `json:"exclude_volumes"`
	IncludeNames   []string `json:"include_names"`
	IncludeVolumes []string `json:"include_volumes"`
}

// Restore defines model for Restore.
type Restore struct {
	Error   *string   `json:"error,omitempty"`
	Id      string    `json:"id"`
	Started time.Time `json:"started"`
	Status  Status    `json:"status"`
	Volumes []Volume  `json:"volumes"`
}

// Rotate defines model for Rotate.
type Rotate struct {
	Error   *string   `json:"error,omitempty"`
	Started time.Time `json:"started"`
	Status  Status    `json:"status"`
}

// RotateInput defines model for RotateInput.
type RotateInput struct {
	Destroy      bool   `json:"destroy"`
	PoliciesPath string `json:"policies_path"`
}

// Status defines model for Status.
type Status string

// StopModes defines model for StopModes.
type StopModes string

// TaskInput defines model for TaskInput.
type TaskInput struct {
	Filters     Filters      `json:"filters"`
	LabelPrefix *string      `json:"label_prefix,omitempty"`
	StopModes   *[]StopModes `json:"stop_modes,omitempty"`
}

// Version defines model for Version.
type Version struct {
	Created string `json:"created"`
	Version string `json:"version"`
}

// Volume defines model for Volume.
type Volume struct {
	Created time.Time `json:"created"`
	Error   *string   `json:"error,omitempty"`
	Ext     string    `json:"ext"`
	Mount   string    `json:"mount"`
	Name    string    `json:"name"`
	Size    *int64    `json:"size,omitempty"`
	Status  Status    `json:"status"`
}

// Backups defines model for Backups.
type Backups = []Backup

// NotFound defines model for NotFound.
type NotFound = Error

// StartBackupJSONRequestBody defines body for StartBackup for application/json ContentType.
type StartBackupJSONRequestBody = TaskInput

// StartBackupWithIdJSONRequestBody defines body for StartBackupWithId for application/json ContentType.
type StartBackupWithIdJSONRequestBody = TaskInput

// StartRestoreJSONRequestBody defines body for StartRestore for application/json ContentType.
type StartRestoreJSONRequestBody = TaskInput

// StartRotateJSONRequestBody defines body for StartRotate for application/json ContentType.
type StartRotateJSONRequestBody = RotateInput
