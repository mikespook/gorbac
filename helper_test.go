package gorbac

import (
	"testing"
)

var (

	/*
	   Suppose:

	   The role-a is inheriting from role-b.
	   The role-b is inheriting from role-c, role-d.
	   The role-c is inheriting from role-a.
	   The role-d is inheriting from role-e.
	   The role-e is individual.

	   There is a circle: role-a, role-c, role-e.
	*/
	circleCases = map[string]map[string][]string{
		RA: {
			"parents": {RB},
		},
		RB: {
			"parents": {RC, RD},
		},
		RC: {
			"parents": {RA},
		},
		RD: {
			"parents": {RE},
		},
		RE: {
			"parents": nil,
		},
	}

	/*
	   Suppose:

	   The role-a is inheriting from role-b.
	   The role-b is inheriting from role-c, role-d.
	   The role-c is inheriting from role-a.
	   The role-d is inheriting from role-e.
	   The role-e is inheriting from role-b.

	   There are two circles: 1) role-a, role-b, role-c; 2) role-b, role-d, role-e.
	*/
	circleCases2 = map[string]map[string][]string{
		RA: {
			"parents": {RB},
		},
		RB: {
			"parents": {RC, RD},
		},
		RC: {
			"parents": {RA},
		},
		RD: {
			"parents": {RE},
		},
		RE: {
			"parents": {RB},
		},
	}
)

func TestInherCircle(t *testing.T) {
	rbac := prepareCase(circleCases)

	if err := InherCircle(rbac); err == nil {
		t.Fatal("There should be a circle inheritance.")
	} else {
		t.Log(err)
	}
}

func TestInherCircle2(t *testing.T) {
	rbac := prepareCase(circleCases2)

	if err := InherCircle(rbac); err == nil {
		t.Fatal("There should be a circle inheritance.")
	} else {
		t.Log(err)
	}
}

func TestInherNormal(t *testing.T) {
	rbac := prepareCase(normalCases)

	if err := InherCircle(rbac); err != nil {
		t.Fatal(err)
	}
}

func TestAllGranted(t *testing.T) {
	rbac := prepareCase(normalCases)

	// All roles have PD
	roles := []string{RA, RB, RE}
	if !AllGranted(rbac, roles, &StdPermission{PD}, nil) {
		t.Errorf("All roles(%v) were expected having %s, but they weren't.", roles, PD)
	}

	roles = []string{RA, RB, RC}
	if AllGranted(rbac, roles, &StdPermission{PD}, nil) {
		t.Errorf("Not all roles(%v) were expected having %s, but they were.", roles, PD)
	}
}

func TestAnyGranted(t *testing.T) {
	rbac := prepareCase(normalCases)

	// All roles have PD
	roles := []string{RA, RB, RE}
	if !AnyGranted(rbac, roles, &StdPermission{PD}, nil) {
		t.Errorf("One of roles(%v) was expected having %s, but it wasn't.", roles, PD)
	}

	roles = []string{RB, RC, RE}
	if AnyGranted(rbac, roles, &StdPermission{PA}, nil) {
		t.Errorf("None of roles(%v) were expected having %s, but it was.", roles, PA)
	}

}

func BenchmarkInherCircle(b *testing.B) {
	rbac := prepareCase(circleCases)

	for i := 0; i < b.N; i++ {
		InherCircle(rbac)
	}
}

func BenchmarkInherCircle2(b *testing.B) {
	rbac := prepareCase(circleCases2)

	for i := 0; i < b.N; i++ {
		InherCircle(rbac)
	}
}

func BenchmarkInherNormal(b *testing.B) {
	rbac := prepareCase(normalCases)

	for i := 0; i < b.N; i++ {
		InherCircle(rbac)
	}
}
