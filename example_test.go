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
	rbac := gorbac.New()
	rA := gorbac.NewStdRole("role-a")
	rB := gorbac.NewStdRole("role-b")
	rC := gorbac.NewStdRole("role-c")
	rD := gorbac.NewStdRole("role-d")
	rE := gorbac.NewStdRole("role-e")

	pA := gorbac.NewStdPermission("permission-a")
	pB := gorbac.NewStdPermission("permission-b")
	pC := gorbac.NewStdPermission("permission-c")
	pD := gorbac.NewStdPermission("permission-d")
	pE := gorbac.NewStdPermission("permission-e")

	rA.AddPermission(pA)
	rB.AddPermission(pB)
	rC.AddPermission(pC)
	rD.AddPermission(pD)
	rE.AddPermission(pE)

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
	// When a circle inheratance ocurred,
	rbac.SetParent("role-c", "role-a")
	// it could be detected as following code:
	if err := gorbac.InherCircle(rbac); err != nil {
		fmt.Println("A circle inheratance ocurred.")
	}
	// Output:
	// The role-a has been granted permis-a, b, c and d.
	// The role-b has been granted permis-b, c and d.
	// A circle inheratance ocurred.
}
