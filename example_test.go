package gorbac_test

import (
	"fmt"

	"github.com/mikespook/gorbac"
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

	The role-a is inheriting from role-b.
	The role-b is inheriting from role-c, role-d.
	The role-c is individual.
	The role-d is individual.
	The role-e is inheriting from role-d.
	Every roles have thire own permissions.
*/
func ExampleRbac() {
	normalCases := map[string]map[string][]string{
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
	rbac := gorbac.New()

	for role, c := range normalCases {
		rbac.Add(role, convPermissions(c["permissions"]), c["parents"])
	}

	if rbac.IsGranted(RA, gorbac.NewStdPermission(PA), nil) &&
		rbac.IsGranted(RA, gorbac.NewStdPermission(PB), nil) &&
		rbac.IsGranted(RA, gorbac.NewStdPermission(PC), nil) &&
		rbac.IsGranted(RA, gorbac.NewStdPermission(PD), nil) {
		fmt.Println("The role-a has been granted permis-a, b, c and d.")
	}
	if rbac.IsGranted(RB, gorbac.NewStdPermission(PB), nil) &&
		rbac.IsGranted(RB, gorbac.NewStdPermission(PC), nil) &&
		rbac.IsGranted(RB, gorbac.NewStdPermission(PD), nil) {
		fmt.Println("The role-b has been granted permis-b, c and d.")
	}
	// When a circle inheratance ocurred,
	rbac.Get(RC).AddParent(RA)
	// it could be detected as following code:
	if err := gorbac.InherCircle(rbac); err != nil {
		fmt.Println("A circle inheratance ocurred.")
	}
	// Output:
	// The role-a has been granted permis-a, b, c and d.
	// The role-b has been granted permis-b, c and d.
	// A circle inheratance ocurred.
}

func convPermissions(a []string) (ps []gorbac.Permission) {
	for _, v := range a {
		ps = append(ps, gorbac.NewStdPermission(v))
	}
	return
}
