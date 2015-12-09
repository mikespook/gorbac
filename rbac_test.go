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
		RA: {
			"permissions": {PA},
			"parents":     {RB},
		},
		RB: {
			"permissions": {PB},
			"parents":     {RC, RD},
		},
		RC: {
			"permissions": {PC},
			"parents":     nil,
		},
		RD: {
			"permissions": {PD},
			"parents":     nil,
		},
		RE: {
			"permissions": nil,
			"parents":     {RD},
		},
	}
)

func convPermissions(a []string) (ps []Permission) {
	for _, v := range a {
		ps = append(ps, &StdPermission{v})
	}
	return
}

func prepareCase(cases map[string]map[string][]string) *RBAC {
	rbac := New()

	for role, c := range cases {
		rbac.Set(role, convPermissions(c["permissions"]), c["parents"])
	}
	return rbac
}

func _TestDumpRestore(t *testing.T) {
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
	if !rbac.IsGranted(RA, &StdPermission{PA}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RA, &StdPermission{PA})
	}
	if !rbac.IsGranted(RA, &StdPermission{PB}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RA, &StdPermission{PB})
	}
	if !rbac.IsGranted(RA, &StdPermission{PC}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RA, &StdPermission{PC})
	}
	if !rbac.IsGranted(RA, &StdPermission{PD}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RA, &StdPermission{PD})
	}
	if rbac.IsGranted(RA, &StdPermission{PE}, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RA, &StdPermission{PE})
	}
}

func TestRbacRoleB(t *testing.T) {
	rbac := prepareCase(normalCases)
	if rbac.IsGranted(RB, &StdPermission{PA}, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RB, &StdPermission{PA})
	}
	if !rbac.IsGranted(RB, &StdPermission{PB}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RB, &StdPermission{PB})
	}
	if !rbac.IsGranted(RB, &StdPermission{PC}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RB, &StdPermission{PC})
	}
	if !rbac.IsGranted(RB, &StdPermission{PD}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RB, &StdPermission{PD})
	}
	if rbac.IsGranted(RB, &StdPermission{PE}, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RB, &StdPermission{PE})
	}
}

func TestRbacRemove(t *testing.T) {
	rbac := prepareCase(normalCases)
	rbac.Remove(RD)
	if rbac.IsGranted(RE, &StdPermission{PD}, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RE, &StdPermission{PD})
	}
	rbac.Get(RA).RevokePermission(&StdPermission{PA})
	if rbac.IsGranted(RA, &StdPermission{PA}, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RA, &StdPermission{PA})
	}
	rbac.Get(RB).RemoveParent(RC)
	if rbac.IsGranted(RB, &StdPermission{PC}, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RB, &StdPermission{PC})
	}
	if !rbac.IsGranted(RB, &StdPermission{PB}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RB, &StdPermission{PB})
	}
	if !rbac.IsGranted(RC, &StdPermission{PC}, nil) {
		t.Errorf("`%s` should be granted `%s`.", RC, &StdPermission{PC})
	}
}

func TestRbacNotExists(t *testing.T) {
	rbac := prepareCase(normalCases)
	if rbac.IsGranted(NOTEXISTS, &StdPermission{PA}, nil) {
		t.Errorf("`%s` should not be granted `%s`.", NOTEXISTS, PA)
	}
	if rbac.IsGranted(RA, &StdPermission{NOTEXISTS}, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RA, NOTEXISTS)
	}
}

func BenchmarkRbacIsGranted(b *testing.B) {
	rbac := prepareCase(normalCases)
	for i := 0; i < b.N; i++ {
		rbac.IsGranted(RE, &StdPermission{PA}, nil)
	}
}

func BenchmarkRbacSet(b *testing.B) {
	rbac := New()
	for i := 0; i < b.N; i++ {
		for role, testingcase := range normalCases {
			rbac.Set(role, convPermissions(testingcase["permissions"]), testingcase["parents"])
		}
	}
}

func BenchmarkRbacAdd(b *testing.B) {
	rbac := New()
	for i := 0; i < b.N; i++ {
		for role, testingcase := range normalCases {
			rbac.Add(role, convPermissions(testingcase["permissions"]), testingcase["parents"])
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
