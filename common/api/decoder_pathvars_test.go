package api_test

import (
	"testing"

	"github.com/janoszen/openshiftci_inspector/common/api"
)

type test struct {
	Field1 string `path:"field1"`
	Field2 int    `path:"field2"`
	Field3 test2
	Field4 test3
}

type test2 struct {
	Field3 bool `path:"field3"`
}

type test3 struct {
	Field4 uint `path:"field4"`
}

func TestDecoding(t *testing.T) {
	decoder := api.NewPathVarsDecoder()
	val := &test{}
	err := decoder.Decode(map[string]string{
		"field1": "foo",
		"field2": "42",
		"field3": "yes",
		"field4": "41",
	}, nil, val)
	if err != nil {
		t.Fatal(err)
	}
	if val.Field1 != "foo" {
		t.Fatal("Failed to set field 1.")
	}
	if val.Field2 != 42 {
		t.Fatal("Failed to set field 2.")
	}
	if val.Field3.Field3 != true {
		t.Fatal("Failed to set field 3.")
	}
	if val.Field4.Field4 != 41 {
		t.Fatal("Failed to set field 3.")
	}
}
