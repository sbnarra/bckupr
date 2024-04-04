package concurrent

import (
	"errors"
	"testing"

	"github.com/sbnarra/bckupr/internal/tests"
	"github.com/sbnarra/bckupr/internal/utils/contexts"
)

func TestRunExecutes(t *testing.T) {
	c := New(tests.Context, "test", 1)
	completed := false

	c.Run(func(ctx contexts.Context) error {
		completed = true
		return nil
	})

	if err := c.Wait(); err != nil {
		t.Fatalf("error waiting for concurrent tasks to complete: %v", err)
	}
	if completed == false {
		t.Fatalf("expected completed bool to be true")
	}
}

func TestMultipleRunExecutes(t *testing.T) {
	c := New(tests.Context, "test", 1)

	completed1 := false
	c.Run(func(ctx contexts.Context) error {
		completed1 = true
		return nil
	})

	completed2 := false
	c.Run(func(ctx contexts.Context) error {
		completed2 = true
		return nil
	})

	if err := c.Wait(); err != nil {
		t.Fatalf("error waiting for concurrent tasks to complete: %v", err)
	}
	if completed1 == false {
		t.Fatalf("expected completed1 bool to be true")
	}
	if completed2 == false {
		t.Fatalf("expected completed2 bool to be true")
	}
}

func TestRunError(t *testing.T) {
	c := New(tests.Context, "test", 1)
	c.Run(func(ctx contexts.Context) error {
		return errors.New("testing")
	})

	if err := c.Wait(); err == nil {
		t.Fatalf("expected error from Wait()")
	}
}
