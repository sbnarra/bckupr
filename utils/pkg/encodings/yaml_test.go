package encodings

import (
	"strings"
	"testing"
)

func TestFromYaml(t *testing.T) {
	data := struct {
		Example string `json:"example"`
	}{}
	err := FromYaml(strings.NewReader(`example: aValue`), &data)

	if data.Example != "aValue" || err != nil {
		t.Fatalf(`wanted data.Example == aValue: %v: %+v`, data.Example, err)
	}
}
