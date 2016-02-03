package gorbac

import (
	"testing"
)

var (
	rA = NewStdRole("role-a")
	pA = NewStdPermission("permission-a")
	rB = NewStdRole("role-b")
	pB = NewStdPermission("permission-b")
	rC = NewStdRole("role-c")
	pC = NewStdPermission("permission-c")

	rbac *RBAC
)

func assert(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestRbacPrepare(t *testing.T) {
	rbac = New()
	assert(t, rA.AddPermission(pA))
	assert(t, rB.AddPermission(pB))
	assert(t, rC.AddPermission(pC))
}

func TestRbacAdd(t *testing.T) {
	assert(t, rbac.Add(rA))
	if err := rbac.Add(rA); err != ErrRoleExist {
		t.Error("A role can not be readded")
	}
	assert(t, rbac.Add(rB))
	assert(t, rbac.Add(rC))
}

func TestRbacRemove(t *testing.T) {
	assert(t, rbac.Remove("role-a"))
	if _, ok := rbac.roles["role-a"]; ok {
		t.Fatal("Role removing failed")
	}
	if err := rbac.Remove("not-exist"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
}

func TestRbacParents(t *testing.T) {
	assert(t, rbac.SetParent("role-c", "role-b"))
	if _, ok := rbac.parents["role-c"]["role-b"]; !ok {
		t.Fatal("Parent binding failed")
	}
	assert(t, rbac.RemoveParent("role-c", "role-b"))
	if _, ok := rbac.parents["role-c"]["role-b"]; ok {
		t.Fatal("Parent unbinding failed")
	}
	if err := rbac.SetParent("role-a", "role-b"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := rbac.SetParent("role-c", "role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	assert(t, rbac.SetParents("role-c", []string{"role-b"}))
	if _, ok := rbac.parents["role-c"]["role-b"]; !ok {
		t.Fatal("Parent binding failed")
	}
}

func TestRbacPermission(t *testing.T) {
	if !rbac.IsGranted("role-c", pC, nil) {
		t.Fatalf("role-c should have %s", pC)
	}
	if rbac.IsGranted("role-c", pC, func(*RBAC, string, Permission) bool { return false }) {
		t.Fatal("Assertion don't work")
	}
	if !rbac.IsGranted("role-c", pB, nil) {
		t.Fatalf("role-c should have %s which inherits from role-b", pB)
	}

	assert(t, rbac.RemoveParent("role-c", "role-b"))
	if rbac.IsGranted("role-c", pB, nil) {
		t.Fatalf("role-c should not have %s because of the unbinding with role-b", pB)
	}
}
