package gorbac

import (
	"testing"
)

func TestStdrA(t *testing.T) {
	rA := NewStdRole("role-a")
	if rA.Id() != "role-a" {
		t.Fatalf("[a] expected, but %s got", rA.Id())
	}
	if err := rA.AssignPermission(NewStdPermission("permission-a")); err != nil {
		t.Fatal(err)
	}
	if !rA.HasPermission(NewStdPermission("permission-a")) {
		t.Fatal("[permission-a] should permit to rA")
	}
	if len(rA.Permissions()) != 1 {
		t.Fatal("[a] should have one permission")
	}

	if err := rA.RevokePermission(NewStdPermission("permission-a")); err != nil {
		t.Fatal(err)
	}
	if rA.HasPermission(NewStdPermission("permission-a")) {
		t.Fatal("[permission-a] should not permit to rA")
	}
	if len(rA.Permissions()) != 0 {
		t.Fatal("[a] should not have any permission")
	}
}
