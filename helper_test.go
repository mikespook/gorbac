package gorbac

import (
	"testing"
)

var (
	pAll  = NewStdPermission("permission-all")
	pNone = NewStdPermission("permission-none")
)

func TestPrepareCircle(t *testing.T) {
	rbac = New()
	assert(t, rA.AddPermission(pA))
	assert(t, rB.AddPermission(pB))
	assert(t, rC.AddPermission(pC))
	assert(t, rA.AddPermission(pAll))
	assert(t, rB.AddPermission(pAll))
	assert(t, rC.AddPermission(pAll))
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

func BenchmarkInherCircle(b *testing.B) {
	rbac = New()
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
	rbac = New()
	rbac.Add(rA)
	rbac.Add(rB)
	rbac.Add(rC)
	rbac.SetParent("role-a", "role-b")
	rbac.SetParent("role-b", "role-c")
	for i := 0; i < b.N; i++ {
		InherCircle(rbac)
	}
}
