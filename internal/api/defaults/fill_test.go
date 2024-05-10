package defaults_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/api/defaults"
	"github.com/sbnarra/bckupr/internal/api/spec"
	"github.com/sbnarra/bckupr/internal/utils/encodings"
)

func TestXxx(t *testing.T) {
	d, err := defaults.New()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	d.AddCreator("time.Time", func() any {
		t := time.Now()
		return &t
	})
	// d.AddCreator("spec.Backup", func() any { return &spec.Backup{} })

	d.AddTypeMapping("spec.Status", func(a any) any {
		return spec.Status(fmt.Sprintf("%v", a))
	})

	v := spec.Backup{}
	fmt.Println(encodings.ToJsonIE(v))
	if err := defaults.Fill(&v, d); err != nil {
		t.Fatalf("%+v", err)
	}
	fmt.Println(encodings.ToJsonIE(v))
}
