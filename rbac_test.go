package gorbac

import (
	"testing"
)

const (
	// Roles
	RA = "role-a"
	RB = "role-b"
	RC = "role-c"
	RD = "role-d"
	RE = "role-e"
	// Permissions
	PA = "permis-a"
	PB = "permis-b"
	PC = "permis-c"
	PD = "permis-d"
	PE = "permis-e"
	//
	NOTEXISTS = "not-exists"
)

var (
	/*
		Suppose:

		The role-a is inheriting from role-b.
		The role-b is inheriting from role-c, role-d.
		The role-c is individual.
		The role-d is individual.
		The role-e is inheriting from role-d.
		Every roles have thire own permissions.
	*/
	normalCases = map[string]map[string][]string{
		RA: map[string][]string{
			"permissions": []string{PA},
			"parents":     []string{RB},
		},
		RB: map[string][]string{
			"permissions": []string{PB},
			"parents":     []string{RC, RD},
		},
		RC: map[string][]string{
			"permissions": []string{PC},
			"parents":     nil,
		},
		RD: map[string][]string{
			"permissions": []string{PD},
			"parents":     nil,
		},
		RE: map[string][]string{
			"permissions": nil,
			"parents":     []string{RD},
		},
	}
)

func prepareCase(cases map[string]map[string][]string) *Rbac {
	rbac := New()

	for role, c := range cases {
		rbac.Set(role, c["permissions"], c["parents"])
	}
	return rbac
}

func TestDumpRestore(t *testing.T) {
	rbac := prepareCase(normalCases)
	m := rbac.Dump()
	a := len(m)
	b := len(rbac.roles)
	if a != b {
		t.Errorf("`%d` roles, dumped `%d`", b, a)
	}
	t.Log(m)
	rbac = Restore(m)
	c := len(rbac.roles)
	if a != c {
		t.Errorf("`%d` roles, restored `%d`", a, c)
	}
}

func TestRbacRoleA(t *testing.T) {
	rbac := prepareCase(normalCases)
	if !rbac.IsGranted(RA, PA, nil) {
		t.Errorf("`%s` should be granted `%s`.", RA, PA)
	}
	if !rbac.IsGranted(RA, PB, nil) {
		t.Errorf("`%s` should be granted `%s`.", RA, PB)
	}
	if !rbac.IsGranted(RA, PC, nil) {
		t.Errorf("`%s` should be granted `%s`.", RA, PC)
	}
	if !rbac.IsGranted(RA, PD, nil) {
		t.Errorf("`%s` should be granted `%s`.", RA, PD)
	}
	if rbac.IsGranted(RA, PE, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RA, PE)
	}
}

func TestRbacRoleB(t *testing.T) {
	rbac := prepareCase(normalCases)
	if rbac.IsGranted(RB, PA, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RB, PA)
	}
	if !rbac.IsGranted(RB, PB, nil) {
		t.Errorf("`%s` should be granted `%s`.", RB, PB)
	}
	if !rbac.IsGranted(RB, PC, nil) {
		t.Errorf("`%s` should be granted `%s`.", RB, PC)
	}
	if !rbac.IsGranted(RB, PD, nil) {
		t.Errorf("`%s` should be granted `%s`.", RB, PD)
	}
	if rbac.IsGranted(RB, PE, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RB, PE)
	}
}

func TestRbacRemove(t *testing.T) {
	rbac := prepareCase(normalCases)
	rbac.Remove(RD)
	if rbac.IsGranted(RE, PD, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RE, PD)
	}
	rbac.Get(RA).RevokePermission(PA)
	if rbac.IsGranted(RA, PA, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RA, PA)
	}
	rbac.Get(RB).RemoveParent(RC)
	if rbac.IsGranted(RB, PC, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RB, PC)
	}
	if !rbac.IsGranted(RB, PB, nil) {
		t.Errorf("`%s` should be granted `%s`.", RB, PB)
	}
	if !rbac.IsGranted(RC, PC, nil) {
		t.Errorf("`%s` should be granted `%s`.", RC, PC)
	}
}

func TestRbacNotExists(t *testing.T) {
	rbac := prepareCase(normalCases)
	if rbac.IsGranted(NOTEXISTS, PA, nil) {
		t.Errorf("`%s` should not be granted `%s`.", NOTEXISTS, PA)
	}
	if rbac.IsGranted(RA, NOTEXISTS, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RA, NOTEXISTS)
	}
}

func BenchmarkRbacIsGranted(b *testing.B) {
	rbac := prepareCase(normalCases)
	for i := 0; i < b.N; i++ {
		rbac.IsGranted(RE, PA, nil)
	}
}

func BenchmarkRbacSet(b *testing.B) {
	rbac := New()
	for i := 0; i < b.N; i++ {
		for role, testingcase := range normalCases {
			rbac.Set(role, testingcase["permissions"], testingcase["parents"])
		}
	}
}

func BenchmarkRbacAdd(b *testing.B) {
	rbac := New()
	for i := 0; i < b.N; i++ {
		for role, testingcase := range normalCases {
			rbac.Add(role, testingcase["permissions"], testingcase["parents"])
		}
	}
}

func BenchmarkDumpRestore(b *testing.B) {
	rbac := prepareCase(normalCases)
	for i := 0; i < b.N; i++ {
		m := rbac.Dump()
		rbac = Restore(m)
	}
}
