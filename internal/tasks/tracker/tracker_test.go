package tracker

import (
	"fmt"
	"testing"
)

type example struct {
	val string
}

func Cleanup() {
	fmt.Println("after")
	tracker = map[string]map[string]*process{}
}

func TestTrackerOnlyAllowsOneKeyId(t *testing.T) {
	t.Cleanup(Cleanup)

	e := example{"foobar"}
	Add("key", "id", &e)
	e2, err := Get[example]("key", "id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.val != e2.val {
		t.Fatalf("values don't match: %v != %v", e, e2)
	}
}

func TestTrackerOnlyAllowsOneProcess(t *testing.T) {
	t.Cleanup(Cleanup)

	e := example{}
	c, err := Add("key", "id", &e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = Add("key", "id", &e)
	if err == nil || err.Error() != "key is already running for id" {
		t.Fatalf("unexpected error: %v", err)
	}

	c(err)
	_, err = Add("key", "id", &e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTrackerOnlyAllowsOneProcess2(t *testing.T) {
	t.Cleanup(Cleanup)

	e := example{}
	c, err := Add("key", "id", &e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = Add("key2", "id", &e)
	if err == nil || err.Error() != "process running: key/id" {
		t.Fatalf("unexpected error: %v", err)
	}

	c(err)
	_, err = Add("key", "id", &e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
