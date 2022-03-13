package gorbac

import (
	"testing"
)

func TestStdrA(t *testing.T) {
	rA := NewRole("role-a")
	if rA.ID != "role-a" {
		t.Fatalf("[a] expected, but %s got", rA.ID)
	}
	if err := rA.Assign(NewBasicPermission("permission-a")); err != nil {
		t.Fatal(err)
	}
	if !rA.Permit(NewBasicPermission("permission-a")) {
		t.Fatal("[permission-a] should permit to rA")
	}
	if len(rA.Permissions()) != 1 {
		t.Fatal("[a] should have one permission")
	}

	if err := rA.Revoke(NewBasicPermission("permission-a")); err != nil {
		t.Fatal(err)
	}
	if rA.Permit(NewBasicPermission("permission-a")) {
		t.Fatal("[permission-a] should not permit to rA")
	}
	if len(rA.Permissions()) != 0 {
		t.Fatal("[a] should not have any permission")
	}
}
