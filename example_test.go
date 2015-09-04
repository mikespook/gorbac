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
	rbac := gorbac.New()

	for role, c := range normalCases {
		rbac.Add(role, c["permissions"], c["parents"])
	}

	if rbac.IsGranted(RA, PA, nil) &&
		rbac.IsGranted(RA, PB, nil) &&
		rbac.IsGranted(RA, PC, nil) &&
		rbac.IsGranted(RA, PD, nil) {
		fmt.Println("The role-a has been granted permis-a, b, c and d.")
	}
	if rbac.IsGranted(RB, PB, nil) &&
		rbac.IsGranted(RB, PC, nil) &&
		rbac.IsGranted(RB, PD, nil) {
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
