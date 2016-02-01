package gorbac

import (
	"encoding/json"
	"testing"
)

func TestStdr1(t *testing.T) {
	r1 := NewStdRole("r1")
	if r1.Id() != "r1" {
		t.Fatalf("[a] expected, but %s got", r1.Id())
	}
	if err := r1.AddPermission(NewStdPermission("p1")); err != nil {
		t.Fatal(err)
	}
	if !r1.HasPermission(NewStdPermission("p1")) {
		t.Fatal("[p1] should permit to r1")
	}

	r2 := NewStdRole("r2")
	if err := r2.AddParent(r1); err != nil {
		t.Fatal(err)
	}
	if !r2.HasPermission(NewStdPermission("p1")) {
		t.Fatal("[p1] should permit to r2")
	}

	if data, err := json.Marshal(r2); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%s", data)
	}

	if err := r1.RevokePermission(NewStdPermission("p1")); err != nil {
		t.Fatal(err)
	}
	if r1.HasPermission(NewStdPermission("p1")) {
		t.Fatal("[p1] should not permit to r1")
	}
	if r2.HasPermission(NewStdPermission("p1")) {
		t.Fatal("[p1] should not permit to r2")
	}
}
