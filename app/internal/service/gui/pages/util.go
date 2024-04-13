package pages

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sbnarra/bckupr/internal/cron"
)

func cronData(cron *cron.Cron) Cron {
	entry := cron.I.Entry(cron.Id)
	return Cron{
		Next:     entry.Next,
		Schedule: cron.Schedule,
	}
}

func load(name string) *template.Template {
	base := "app/"
	if val, exists := os.LookupEnv("UI_BASE_PATH"); exists {
		base = val
	}

	t := template.Must(template.ParseFiles(fmt.Sprintf(base+"ui/%v.html", name)))
	t = template.Must(t.ParseGlob(base + "ui/common/*"))
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

	partsPath := fmt.Sprintf("ui/%v", name)
	if _, err := os.Stat(partsPath); err != nil {
		//TODO: delete this noise
		// fmt.Printf("%+v\n", err)
	} else {
		t = template.Must(t.ParseGlob(partsPath + "/*"))
	}

	return t
}
