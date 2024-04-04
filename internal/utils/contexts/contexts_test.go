package contexts_test

import (
	"testing"

	"github.com/sbnarra/bckupr/internal/utils/contexts"
	"github.com/sbnarra/bckupr/internal/utils/logging"
)

var Test = contexts.Create("under-test", "/tmp/backuprs", false, true, func(ctx contexts.Context, a any) {
	logging.Info(ctx, a)
})

func TestFeedback(t *testing.T) {
	var data any
	context := contexts.Create("under-test", "/tmp/backuprs", false, true, func(ctx contexts.Context, a any) {
		data = a
	})

	context.Feedback("hello")
	if data != "hello" {
		t.Errorf("expect data to be 'hello': %v", data)
	}
}

func TestFeedbackData(t *testing.T) {
	var data any
	context := contexts.Create("under-test", "/tmp/backuprs", false, true, func(ctx contexts.Context, a any) {
		data = a
	})

	context.FeedbackJson("hello")
	if data != "\"hello\"" {
		t.Errorf("expect data to be '\"hello\"': %v", data)
	}
}
