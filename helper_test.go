package gorbac

import (
	"errors"
	"testing"
)

var (
	pAll  = NewPermission("permission-all")
	pNone = NewPermission("permission-none")
)

func TestPrepareCircle(t *testing.T) {
	rbac = New[string]()
	assert(t, rA.Assign(pA))
	assert(t, rB.Assign(pB))
	assert(t, rC.Assign(pC))
	assert(t, rA.Assign(pAll))
	assert(t, rB.Assign(pAll))
	assert(t, rC.Assign(pAll))
	assert(t, rbac.Add(rA))
	assert(t, rbac.Add(rB))
	assert(t, rbac.Add(rC))
	assert(t, rbac.SetParent("role-a", "role-b"))
	assert(t, rbac.SetParent("role-b", "role-c"))
	assert(t, rbac.SetParent("role-c", "role-a"))
}

func TestInherCircle(t *testing.T) {
	if err := InherCircle(rbac); err == nil {
		t.Fatal("There should be a circle inheritance.")
	} else {
		t.Log(err)
	}
}

func TestInherNormal(t *testing.T) {
	assert(t, rbac.RemoveParent("role-c", "role-a"))
	if err := InherCircle(rbac); err != nil {
		t.Fatal(err)
	}
}

func TestAllGranted(t *testing.T) {
	// All roles have pAll
	roles := []string{"role-a", "role-b", "role-c"}
	if !AllGranted(rbac, roles, pAll, nil) {
		t.Errorf("All roles(%v) were expected having %s, but they weren't.", roles, pAll)
	}

	if AllGranted(rbac, roles, pA, nil) {
		t.Errorf("Not all roles(%v) were expected having %s, but they were.", roles, pA)
	}
}

func TestAnyGranted(t *testing.T) {
	// rA roles have pA
	roles := []string{"role-a", "role-b", "role-c"}
	if !AnyGranted(rbac, roles, pA, nil) {
		t.Errorf("One of roles(%v) was expected having %s, but it wasn't.", roles, pA)
	}

	if AnyGranted(rbac, roles, pNone, nil) {
		t.Errorf("None of roles(%v) were expected having %s, but it was.", roles, pNone)
	}

}

func TestWalk(t *testing.T) {
	if err := Walk(rbac, nil); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	h := func(r Role[string], parents []string) error {
		t.Logf("Role: %v", r.ID)
		permissions := make([]string, 0)
		for _, p := range r.Permissions() {
			permissions = append(permissions, p.ID())
		}
		t.Logf("Permission: %v", permissions)
		t.Logf("Parents: %v", parents)
		return nil
	}
	if err := Walk(rbac, h); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	he := func(r Role[string], parents []string) error {
		return errors.New("Expected error")
	}
	if err := Walk(rbac, he); err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func BenchmarkInherCircle(b *testing.B) {
	rbac = New[string]()
	rbac.Add(rA)
	rbac.Add(rB)
	rbac.Add(rC)
	rbac.SetParent("role-a", "role-b")
	rbac.SetParent("role-b", "role-c")
	rbac.SetParent("role-c", "role-a")
	for i := 0; i < b.N; i++ {
		InherCircle(rbac)
	}
}

func BenchmarkInherNormal(b *testing.B) {
	rbac = New[string]()
	rbac.Add(rA)
	rbac.Add(rB)
	rbac.Add(rC)
	rbac.SetParent("role-a", "role-b")
	rbac.SetParent("role-b", "role-c")
	for i := 0; i < b.N; i++ {
		InherCircle(rbac)
	}
}
