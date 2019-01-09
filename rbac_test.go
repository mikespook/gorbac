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
	assert(t, rA.Assign(pA))
	assert(t, rB.Assign(pB))
	assert(t, rC.Assign(pC))
}

func TestRbacAdd(t *testing.T) {
	assert(t, rbac.Add(rA))
	if err := rbac.Add(rA); err != ErrRoleExist {
		t.Error("A role can not be readded")
	}
	assert(t, rbac.Add(rB))
	assert(t, rbac.Add(rC))
}

func TestRbacGetRemove(t *testing.T) {
	assert(t, rbac.SetParent("role-c", "role-a"))
	assert(t, rbac.SetParent("role-a", "role-b"))
	if r, parents, err := rbac.Get("role-a"); err != nil {
		t.Fatal(err)
	} else if r.ID() != "role-a" {
		t.Fatalf("[role-a] does not match %s", r.ID())
	} else if len(parents) != 1 {
		t.Fatal("[role-a] should have one parent")
	}
	assert(t, rbac.Remove("role-a"))
	if _, ok := rbac.roles["role-a"]; ok {
		t.Fatal("Role removing failed")
	}
	if err := rbac.Remove("not-exist"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if r, parents, err := rbac.Get("role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	} else if r != nil {
		t.Fatal("The instance of role should be a nil")
	} else if parents != nil {
		t.Fatal("The slice of parents should be a nil")
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
	if err := rbac.RemoveParent("role-a", "role-b"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := rbac.RemoveParent("role-b", "role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := rbac.SetParent("role-a", "role-b"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := rbac.SetParent("role-c", "role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := rbac.SetParents("role-a", []string{"role-b"}); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	if err := rbac.SetParents("role-c", []string{"role-a"}); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	}
	assert(t, rbac.SetParents("role-c", []string{"role-b"}))
	if _, ok := rbac.parents["role-c"]["role-b"]; !ok {
		t.Fatal("Parent binding failed")
	}
	if parents, err := rbac.GetParents("role-a"); err != ErrRoleNotExist {
		t.Fatalf("%s needed", ErrRoleNotExist)
	} else if len(parents) != 0 {
		t.Fatal("[role-a] should not have any parent")
	}
	if parents, err := rbac.GetParents("role-b"); err != nil {
		t.Fatal(err)
	} else if len(parents) != 0 {
		t.Fatal("[role-b] should not have any parent")
	}
	if parents, err := rbac.GetParents("role-c"); err != nil {
		t.Fatal(err)
	} else if len(parents) != 1 {
		t.Fatal("[role-c] should have one parent")
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

	if rbac.IsGranted("role-a", nil, nil) {
		t.Fatal("role-a should not have nil permission")
	}
}

func BenchmarkRbacGranted(b *testing.B) {
	rbac = New()
	rA.Assign(pA)
	rB.Assign(pB)
	rC.Assign(pC)
	rbac.Add(rA)
	rbac.Add(rB)
	rbac.Add(rC)
	for i := 0; i < b.N; i++ {
		rbac.IsGranted("role-a", pA, nil)
	}
}

func BenchmarkRbacNotGranted(b *testing.B) {
	rbac = New()
	rA.Assign(pA)
	rB.Assign(pB)
	rC.Assign(pC)
	rbac.Add(rA)
	rbac.Add(rB)
	rbac.Add(rC)
	for i := 0; i < b.N; i++ {
		rbac.IsGranted("role-a", pB, nil)
	}
}
