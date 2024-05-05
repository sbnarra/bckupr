package pages

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sbnarra/bckupr/internal/cron"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/errors"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

func cronData(cron *cron.Cron) Cron {
	backup := cron.I.Entry(cron.BackupId)
	backupSchedule := cron.BackupSchedule
	nextBackup := backup.Next.Format("2006-01-02 15:04:00 MST")
	if cron.BackupSchedule == "" {
		backupSchedule = "disabled"
		nextBackup = "disabled"
	}

	rotate := cron.I.Entry(cron.RotateId)
	rotateSchedule := cron.RotateSchedule
	nextRotate := rotate.Next.Format("2006-01-02 15:04:00 UTC")
	if cron.RotateSchedule == "" {
		rotateSchedule = "disabled"
		nextRotate = "disabled"
	}

	return Cron{
		NextBackup:     nextBackup,
		BackupSchedule: backupSchedule,
		NextRotate:     nextRotate,
		RotateSchedule: rotateSchedule,
	}
}

func loadAndExecute(ctx contexts.Context, name string, wr io.Writer, data any) *errors.Error {
	loaded := load(ctx, name)
	err := loaded.Execute(wr, data)
	return errors.Wrap(err, "error executing template "+name)
}

func load(ctx contexts.Context, name string) *template.Template {
	base := ""
	if val, exists := os.LookupEnv("UI_BASE_PATH"); exists {
		base = val
	}
	page := fmt.Sprintf(base+"web/%v.html", name)
	logging.Debug(ctx, "loading page", page)

	t := template.Must(template.ParseFiles(page))
	t = template.Must(t.ParseGlob(base + "web/common/*"))
	t = t.Funcs(template.FuncMap{
		"date": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"time": func(t time.Time) string {
			return t.Format("15:04:05 MST")
		},
		"datetime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05 MST")
		},
		"noFileExtension": func(filename string) string {
			filename = strings.TrimSuffix(filename, filepath.Ext(filename))
			filename = strings.TrimSuffix(filename, filepath.Ext(filename))
			return filename
		},
		"map": func(items ...any) map[string]any {
			m := map[string]any{}
			for i := 0; i < len(items); i = i + 2 {
				m[fmt.Sprintf("%v", items[i])] = items[i+1]
			}
			return m
		},
	})

	if st, err := os.Stat(fmt.Sprintf("web/%v", name)); err == nil && st.IsDir() {
		t = template.Must(t.ParseGlob(fmt.Sprintf("web/%v/*", name)))
	} else {
		// logging.CheckError(ctx, err)
	}

	return t
}
