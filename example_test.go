package gorbac_test

import (
	"fmt"

	"github.com/mikespook/gorbac"
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
	rbac := gorbac.New[string]()
	rA := gorbac.NewRole("role-a")
	rB := gorbac.NewRole("role-b")
	rC := gorbac.NewRole("role-c")
	rD := gorbac.NewRole("role-d")
	rE := gorbac.NewRole("role-e")

	pA := gorbac.NewBasicPermission("permission-a")
	pB := gorbac.NewBasicPermission("permission-b")
	pC := gorbac.NewBasicPermission("permission-c")
	pD := gorbac.NewBasicPermission("permission-d")
	pE := gorbac.NewBasicPermission("permission-e")

	rA.Assign(pA)
	rB.Assign(pB)
	rC.Assign(pC)
	rD.Assign(pD)
	rE.Assign(pE)

	rbac.Add(rA)
	rbac.Add(rB)
	rbac.Add(rC)
	rbac.Add(rD)
	rbac.Add(rE)
	rbac.SetParent("role-a", "role-b")
	rbac.SetParents("role-b", []string{"role-c", "role-d"})
	rbac.SetParent("role-e", "role-d")

	if rbac.IsGranted("role-a", pA, nil) &&
		rbac.IsGranted("role-a", pB, nil) &&
		rbac.IsGranted("role-a", pC, nil) &&
		rbac.IsGranted("role-a", pD, nil) {
		fmt.Println("The role-a has been granted permis-a, b, c and d.")
	}
	if rbac.IsGranted("role-b", pB, nil) &&
		rbac.IsGranted("role-b", pC, nil) &&
		rbac.IsGranted("role-b", pD, nil) {
		fmt.Println("The role-b has been granted permis-b, c and d.")
	}
	// When a circle inheratance occurred,
	rbac.SetParent("role-c", "role-a")
	// it could be detected as following code:
	if err := gorbac.InherCircle(rbac); err != nil {
		fmt.Println("A circle inheratance occurred.")
	}
	// Output:
	// The role-a has been granted permis-a, b, c and d.
	// The role-b has been granted permis-b, c and d.
	// A circle inheratance occurred.
}
