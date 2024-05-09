// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package spec

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
)

// Defines values for StopModes.
const (
	All      StopModes = "all"
	Attached StopModes = "attached"
	Labelled StopModes = "labelled"
	Linked   StopModes = "linked"
	Writers  StopModes = "writers"
)

// Defines values for TaskStatus.
const (
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusError     TaskStatus = "error"
	TaskStatusPending   TaskStatus = "pending"
)

// Backup defines model for Backup.
type Backup struct {
	Created time.Time `json:"created"`
	Id      string    `json:"id"`
	Type    string    `json:"type"`
	Volumes []Volume  `json:"volumes"`
	union   json.RawMessage
}

// BackupTrigger defines model for BackupTrigger.
type BackupTrigger struct {
	union json.RawMessage
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

// RestoreTrigger defines model for RestoreTrigger.
type RestoreTrigger struct {
	Dummy *string `json:"dummy,omitempty"`
	union json.RawMessage
}

// RotateTrigger defines model for RotateTrigger.
type RotateTrigger struct {
	Destroy      bool   `json:"destroy"`
	PoliciesPath string `json:"policies_path"`
}

// StopModes defines model for StopModes.
type StopModes string

// Task defines model for Task.
type Task struct {
	Created time.Time  `json:"created"`
	Status  TaskStatus `json:"status"`
}

// TaskStatus defines model for Task.Status.
type TaskStatus string

// TaskTrigger defines model for TaskTrigger.
type TaskTrigger struct {
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
	Error   string    `json:"error"`
	Ext     string    `json:"ext"`
	Mount   string    `json:"mount"`
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
}

// Backups defines model for Backups.
type Backups = []Backup

// NotFound defines model for NotFound.
type NotFound = Error

// TriggerBackupJSONRequestBody defines body for TriggerBackup for application/json ContentType.
type TriggerBackupJSONRequestBody = BackupTrigger

// TriggerBackupWithIdJSONRequestBody defines body for TriggerBackupWithId for application/json ContentType.
type TriggerBackupWithIdJSONRequestBody = BackupTrigger

// TriggerRestoreJSONRequestBody defines body for TriggerRestore for application/json ContentType.
type TriggerRestoreJSONRequestBody = RestoreTrigger

// RotateBackupsJSONRequestBody defines body for RotateBackups for application/json ContentType.
type RotateBackupsJSONRequestBody = RotateTrigger

// AsTask returns the union data inside the Backup as a Task
func (t Backup) AsTask() (Task, error) {
	var body Task
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromTask overwrites any union data inside the Backup as the provided Task
func (t *Backup) FromTask(v Task) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeTask performs a merge with any union data inside the Backup, using the provided Task
func (t *Backup) MergeTask(v Task) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t Backup) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	if err != nil {
		return nil, err
	}
	object := make(map[string]json.RawMessage)
	if t.union != nil {
		err = json.Unmarshal(b, &object)
		if err != nil {
			return nil, err
		}
	}

	object["created"], err = json.Marshal(t.Created)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'created': %w", err)
	}

	object["id"], err = json.Marshal(t.Id)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'id': %w", err)
	}

	object["type"], err = json.Marshal(t.Type)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'type': %w", err)
	}

	object["volumes"], err = json.Marshal(t.Volumes)
	if err != nil {
		return nil, fmt.Errorf("error marshaling 'volumes': %w", err)
	}

	b, err = json.Marshal(object)
	return b, err
}

func (t *Backup) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	object := make(map[string]json.RawMessage)
	err = json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["created"]; found {
		err = json.Unmarshal(raw, &t.Created)
		if err != nil {
			return fmt.Errorf("error reading 'created': %w", err)
		}
	}

	if raw, found := object["id"]; found {
		err = json.Unmarshal(raw, &t.Id)
		if err != nil {
			return fmt.Errorf("error reading 'id': %w", err)
		}
	}

	if raw, found := object["type"]; found {
		err = json.Unmarshal(raw, &t.Type)
		if err != nil {
			return fmt.Errorf("error reading 'type': %w", err)
		}
	}

	if raw, found := object["volumes"]; found {
		err = json.Unmarshal(raw, &t.Volumes)
		if err != nil {
			return fmt.Errorf("error reading 'volumes': %w", err)
		}
	}

	return err
}

// AsTaskTrigger returns the union data inside the BackupTrigger as a TaskTrigger
func (t BackupTrigger) AsTaskTrigger() (TaskTrigger, error) {
	var body TaskTrigger
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromTaskTrigger overwrites any union data inside the BackupTrigger as the provided TaskTrigger
func (t *BackupTrigger) FromTaskTrigger(v TaskTrigger) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeTaskTrigger performs a merge with any union data inside the BackupTrigger, using the provided TaskTrigger
func (t *BackupTrigger) MergeTaskTrigger(v TaskTrigger) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t BackupTrigger) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *BackupTrigger) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsTaskTrigger returns the union data inside the RestoreTrigger as a TaskTrigger
func (t RestoreTrigger) AsTaskTrigger() (TaskTrigger, error) {
	var body TaskTrigger
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromTaskTrigger overwrites any union data inside the RestoreTrigger as the provided TaskTrigger
func (t *RestoreTrigger) FromTaskTrigger(v TaskTrigger) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeTaskTrigger performs a merge with any union data inside the RestoreTrigger, using the provided TaskTrigger
func (t *RestoreTrigger) MergeTaskTrigger(v TaskTrigger) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JSONMerge(t.union, b)
	t.union = merged
	return err
}

func (t RestoreTrigger) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	if err != nil {
		return nil, err
	}
	object := make(map[string]json.RawMessage)
	if t.union != nil {
		err = json.Unmarshal(b, &object)
		if err != nil {
			return nil, err
		}
	}

	if t.Dummy != nil {
		object["dummy"], err = json.Marshal(t.Dummy)
		if err != nil {
			return nil, fmt.Errorf("error marshaling 'dummy': %w", err)
		}
	}
	b, err = json.Marshal(object)
	return b, err
}

func (t *RestoreTrigger) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	object := make(map[string]json.RawMessage)
	err = json.Unmarshal(b, &object)
	if err != nil {
		return err
	}

	if raw, found := object["dummy"]; found {
		err = json.Unmarshal(raw, &t.Dummy)
		if err != nil {
			return fmt.Errorf("error reading 'dummy': %w", err)
		}
	}

	return err
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /backups)
	ListBackups(c *gin.Context)
	// Creates new backup
	// (POST /backups)
	TriggerBackup(c *gin.Context)
	// Deletes backup
	// (DELETE /backups/{id})
	DeleteBackup(c *gin.Context, id string)
	// Gets backup by id
	// (GET /backups/{id})
	GetBackup(c *gin.Context, id string)

	// (PUT /backups/{id})
	TriggerBackupWithId(c *gin.Context, id string)

	// (POST /backups/{id}/restore)
	TriggerRestore(c *gin.Context, id string)
	// Retrieves application version
	// (POST /rotate)
	RotateBackups(c *gin.Context)
	// Retrieves application version
	// (GET /version)
	GetVersion(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// ListBackups operation middleware
func (siw *ServerInterfaceWrapper) ListBackups(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListBackups(c)
}

// TriggerBackup operation middleware
func (siw *ServerInterfaceWrapper) TriggerBackup(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.TriggerBackup(c)
}

// DeleteBackup operation middleware
func (siw *ServerInterfaceWrapper) DeleteBackup(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.DeleteBackup(c, id)
}

// GetBackup operation middleware
func (siw *ServerInterfaceWrapper) GetBackup(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetBackup(c, id)
}

// TriggerBackupWithId operation middleware
func (siw *ServerInterfaceWrapper) TriggerBackupWithId(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.TriggerBackupWithId(c, id)
}

// TriggerRestore operation middleware
func (siw *ServerInterfaceWrapper) TriggerRestore(c *gin.Context) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", c.Param("id"), &id, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter id: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.TriggerRestore(c, id)
}

// RotateBackups operation middleware
func (siw *ServerInterfaceWrapper) RotateBackups(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.RotateBackups(c)
}

// GetVersion operation middleware
func (siw *ServerInterfaceWrapper) GetVersion(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetVersion(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/backups", wrapper.ListBackups)
	router.POST(options.BaseURL+"/backups", wrapper.TriggerBackup)
	router.DELETE(options.BaseURL+"/backups/:id", wrapper.DeleteBackup)
	router.GET(options.BaseURL+"/backups/:id", wrapper.GetBackup)
	router.PUT(options.BaseURL+"/backups/:id", wrapper.TriggerBackupWithId)
	router.POST(options.BaseURL+"/backups/:id/restore", wrapper.TriggerRestore)
	router.POST(options.BaseURL+"/rotate", wrapper.RotateBackups)
	router.GET(options.BaseURL+"/version", wrapper.GetVersion)
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xY224bNxD9lQXbRzYSWqMo9i1pk8BAb0gD98EQDGp3JDHmklty1okq7L8XvO1NXEuJ",
	"7KBva17OzJw5nBn5QApV1UqCREPyA9FgaiUNuD9eseK+qe1XoSSCRPvJ6lrwgiFXcvHBKGnXTLGDitmv",
	"bzVsSE6+WfSwC79rFgGubVtKSjCF5rVFIXk01NLwZT7LJkeozLnGKcF9DSQnTGu2n3fG2KO/K3yjGlk+",
	"GQWvtVY6ZfS1RI77TCrMNs5iS8l7Zu6fzLIDSxgO6zTATPLO5P6PDclvz8FeUVJrVYNG7uVTaGAIjryN",
	"0hVDkpOSIXyHvALSJcKg5nJrA+bu7NGyX0hsPCjRVN7WWRq4cecTGqBEwz8N19bZW+sG7ZwPZ3tjq+62",
	"Wn+AAnvVvtd8uwX9ebzFS20K18slP0yIhbg8YWQShz+Wwn3DBYI2CeRPhWhKuJMs8FrChjUCSX67oj3J",
	"MymKfNIOZ5ChL0Ti8mk8ijiXejTVysi9YzN0QukxNan8vAODSsPFghont2yqap+WzbEHChkOHZhAgUGt",
	"hmBrpQQwae/WSvCCg7mrGe5Oy3R8nHbYKWL+QlX/psogVtlUFoAJQSgRbA1CuBf7UXMnb0oYIit2blFw",
	"eQ/lALXPcCy0F5YvgwybkWc1yNJuUtdlBfiCMn2XM7z0FSgApwgZZvwogk3/zB9TTawGbSDxrtaw4Z9G",
	"T4SsbYXT6bBVfVfFpJxViPs0nnpfMYRU7DegDfc9cDZzxz2jv9QHt3yxfLEkpxISr/bNIemVbzKXy2mu",
	"ytsSgsn1SjUyvWOLT3LD8H9h5BKX+ONV7w6XCK6aTLhwgN6TaHfYMx0sne1ArSvIGxVKyXD4cjLLXv55",
	"bX3gKKBbJIPkEcEQjONb1SBZzUlOfghZtHXEUb5Y99PkFhwzNiNufLouSU5+5QbjuEfHs+/3y+Wcirtz",
	"i25UtDJmW2OZ8TbJyhVCkzAanmsYsjytYPCVKvdPPGp3rWBuyM1QZWvI+rz1KUbdQPvlnPiRsqkqpvck",
	"Jz87CyaT8DFbx8CPGGtpl7LFgZetl4ctm8cs/uLWByQe+zmO2I3V0fjYO49lHvGMpvXzFvBRD85jipKr",
	"5dXp492PkbHzbwGj59l6n7nhNaVFplkFvh3cTrkJ191dbhdCK/Zlww/EY2nQgRCnVdMOHs0p4f/NcXdd",
	"/q/kPyP3Z/7tG6xn63QVmb4JKwg7G7oG83Vz+lg1CxPrM+VzMg9/1XrW/UBOpka7OdklI0mPn6OHPeYZ",
	"2BnN6heq/Uwy+vrzDlBzeACTDSLI+lkp0mb2BqEKtA2msLnCetMhPNujjCYSlL0cxNIfuyxuex30Q3yu",
	"jRYkJzvEOl8shCqY2CmD+U/L5dL/N8ADHOKLja2LdisBul21/wUAAP//w/T16scTAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}