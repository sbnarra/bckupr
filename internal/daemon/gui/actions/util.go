package actions

import (
	"net/http"
)

func readForm(r *http.Request) (map[string][]string, error) {
	form := map[string][]string{}
	if err := r.ParseForm(); err != nil {
		return form, err
	} else {
		for key, value := range r.Form {
			form[key] = value
		}
	}
	return form, nil
}
