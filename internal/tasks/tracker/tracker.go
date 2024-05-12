package tracker

import (
	"time"

	"github.com/sbnarra/bckupr/internal/utils/errors"
)

type process struct {
	data     any
	complete bool
	started  time.Time
	finished time.Time
	err      *errors.E
}

var tracker = map[string]map[string]*process{}

func processRunning() (string, string) {
	for key, processes := range tracker {
		for id, process := range processes {
			if !process.complete {
				return key, id
			}
		}
	}
	return "", ""
}

func Add(key string, id string, data any) (func(*errors.E), *errors.E) {
	if processes, found := tracker[key]; !found {
		tracker[key] = map[string]*process{}
	} else if process, exists := processes[id]; exists && !process.complete {
		return nil, errors.Errorf("%v is already running for %v", key, id)
	}

	if key, id := processRunning(); key != "" {
		return nil, errors.Errorf("process running: %v/%v", key, id)
	}

	process := process{
		data:     data,
		started:  time.Now(),
		complete: false,
	}
	close := func(err *errors.E) {
		process.complete = true
		process.finished = time.Now()
		process.err = err
	}
	tracker[key][id] = &process
	return close, nil
}

// for services not persisting data, e.g restore/rotate
func Get[T any](key string, id string) (*T, *errors.E) {
	if processes, found := tracker[key]; !found {
		return nil, errors.Errorf("no processes found for %v", key)
	} else if process, exists := processes[id]; !exists {
		return nil, errors.Errorf("process %v not found for %v", id, key)
	} else {
		data := process.data.(*T)
		return data, nil
	}
}
