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

/*
Suppose:

The role-c is inheriting from role-a, role-b.
The role-d is individual.
The role-e is inheriting from role-c, role-d.
Every roles have thire own permissions.

*/
var (
	testingcases = map[string]map[string][]string {
		RA: map[string][]string{
			"permissions": []string{PA},
			"parents": nil,
		},
		RB: map[string][]string{
			"permissions": []string{PB},
			"parents": nil,
		},
		RC: map[string][]string{
			"permissions": []string{PC},
			"parents": []string{RA, RB},
		},
		RD: map[string][]string{
			"permissions": []string{PD},
			"parents": nil,
		},
		RE: map[string][]string{
			"permissions": nil,
			"parents": []string{RC, RD},
		},
	}

	rbac = New()
)

func init() {
	for role, testingcase := range testingcases {
		rbac.Set(role, testingcase["permissions"], testingcase["parents"])
	}
}

/*
Thus:

The role-c has been granted permis-a/b/c.
The role-e has been granted permis-a/b/c/d.
*/
func TestRbacRoleC(t *testing.T) {
	if !rbac.IsGranted(RC, PA, nil) {
		t.Errorf("`%s` should be granted `%s`.", RC, PA)
	}
	if !rbac.IsGranted(RC, PB, nil) {
		t.Errorf("`%s` should be granted `%s`.", RC, PB)
	}
	if !rbac.IsGranted(RC, PC, nil) {
		t.Errorf("`%s` should be granted `%s`.", RC, PC)
	}
	if rbac.IsGranted(RC, PD, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RC, PD)
	}
}

func TestRbacRoleE(t *testing.T) {
	if !rbac.IsGranted(RE, PA, nil) {
		t.Errorf("`%s` should be granted `%s`.", RE, PA)
	}
	if !rbac.IsGranted(RE, PC, nil) {
		t.Errorf("`%s` should be granted `%s`.", RE, PC)
	}
	if !rbac.IsGranted(RE, PD, nil) {
		t.Errorf("`%s` should be granted `%s`.", RE, PD)
	}
}

func TestRbacRemove(t *testing.T) {
	rbac.Remove(RD)
	if rbac.IsGranted(RE, PD, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RE, PD)
	}
	rbac.GetRole(RA).RevokePermission(PA)
	if rbac.IsGranted(RA, PA, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RA, PA)
	}
	rbac.GetRole(RC).RemoveParent(RB)
	if rbac.IsGranted(RC, PB, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RC, PB)
	}
	if !rbac.IsGranted(RB, PB, nil) {
		t.Errorf("`%s` should be granted `%s`.", RC, PC)
	}
	if !rbac.IsGranted(RC, PC, nil) {
		t.Errorf("`%s` should be granted `%s`.", RC, PC)
	}
}

func TestRbacNotExists(t *testing.T) {
	if rbac.IsGranted(NOTEXISTS, PA, nil) {
		t.Errorf("`%s` should not be granted `%s`.", NOTEXISTS, PA)
	}
	if rbac.IsGranted(RA, NOTEXISTS, nil) {
		t.Errorf("`%s` should not be granted `%s`.", RA, NOTEXISTS)
	}
}

func BenchmarkRbacIsGranted(b *testing.B) {
	rbac := New()
	for role, testingcase := range testingcases {
		rbac.Set(role, testingcase["permissions"], testingcase["parents"])
	}
	for i := 0; i < b.N; i ++ {
		rbac.IsGranted(RE, PA, nil)
	}
}

func BenchmarkRbacSet(b *testing.B) {
	rbac := New()
	for i := 0; i < b.N; i ++ {
		for role, testingcase := range testingcases {
			rbac.Set(role, testingcase["permissions"], testingcase["parents"])
		}
	}
}

func BenchmarkRbacAdd(b *testing.B) {
	rbac := New()
	for i := 0; i < b.N; i ++ {
		for role, testingcase := range testingcases {
			rbac.Add(role, testingcase["permissions"], testingcase["parents"])
		}
	}
}
