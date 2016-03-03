package gorbac

import (
	"testing"
)

func TestStdrA(t *testing.T) {
	rA := NewStdRole("role-a")
	if rA.ID() != "role-a" {
		t.Fatalf("[a] expected, but %s got", rA.ID())
	}
	if err := rA.Assign(NewStdPermission("permission-a")); err != nil {
		t.Fatal(err)
	}
	if !rA.Permit(NewStdPermission("permission-a")) {
		t.Fatal("[permission-a] should permit to rA")
	}
	if len(rA.Permissions()) != 1 {
		t.Fatal("[a] should have one permission")
	}

	if err := rA.Revoke(NewStdPermission("permission-a")); err != nil {
		t.Fatal(err)
	}
	if rA.Permit(NewStdPermission("permission-a")) {
		t.Fatal("[permission-a] should not permit to rA")
	}
	if len(rA.Permissions()) != 0 {
		t.Fatal("[a] should not have any permission")
	}
}
