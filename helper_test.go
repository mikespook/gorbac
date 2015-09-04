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
		RA: map[string][]string{
			"parents": []string{RB},
		},
		RB: map[string][]string{
			"parents": []string{RC, RD},
		},
		RC: map[string][]string{
			"parents": []string{RA},
		},
		RD: map[string][]string{
			"parents": []string{RE},
		},
		RE: map[string][]string{
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
		RA: map[string][]string{
			"parents": []string{RB},
		},
		RB: map[string][]string{
			"parents": []string{RC, RD},
		},
		RC: map[string][]string{
			"parents": []string{RA},
		},
		RD: map[string][]string{
			"parents": []string{RE},
		},
		RE: map[string][]string{
			"parents": []string{RB},
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
