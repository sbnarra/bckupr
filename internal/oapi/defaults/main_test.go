package defaults_test

import (
	"testing"
	"time"

	"github.com/sbnarra/bckupr/internal/oapi/defaults"
)

func TestXxx(t *testing.T) {
	specPath := "../../../api/openapi-system.yml"
	d, err := defaults.New(specPath)
	if err != nil {
		t.Fatalf("failed to load spec from %v: %+v", specPath, err)
	}

	d.AddType("time.Time", func() interface{} {
		t := time.Now()
		return &t
	})
	// d.AddType("client.Other", func() interface{} { return &client.Other{} })

	// v := client.Version{}
	// if err := defaults.Fill(&v, d); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(v)
	// fmt.Println("Created:", v.Created)
	// fmt.Println("OtherField:", v.Other.OtherField)
}
