package encodings

import (
	"testing"
)

func TestToJson(t *testing.T) {
	if json, err := ToJson("hello"); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else {
		if json != "\"hello\"" {
			t.Errorf("expected json encoding: '\"hello\"' but got %v", json)
		}
	}
}
