package gorbac_test

import (
	"fmt"

	"github.com/mikespook/gorbac/v3"
)

/*
	Suppose:

	The role-a is inheriting from role-b.
	The role-b is inheriting from role-c, role-d.
	The role-c is individual.
	The role-d is individual.
	The role-e is inheriting from role-d.
	Every roles have their own permissions.
*/
func ExampleRbacStringID() {
	rbac := gorbac.New[string]()
	rA := gorbac.NewRole("role-a")
	rB := gorbac.NewRole("role-b")
	rC := gorbac.NewRole("role-c")
	rD := gorbac.NewRole("role-d")
	rE := gorbac.NewRole("role-e")

	pA := gorbac.NewPermission("permission-a")
	pB := gorbac.NewPermission("permission-b")
	pC := gorbac.NewPermission("permission-c")
	pD := gorbac.NewPermission("permission-d")
	pE := gorbac.NewPermission("permission-e")

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

func ExampleRbacNumberID() {
	rbac := gorbac.New[int]()
	rA := gorbac.NewRole(1)
	rB := gorbac.NewRole(2)
	rC := gorbac.NewRole(3)
	rD := gorbac.NewRole(4)
	rE := gorbac.NewRole(5)

	pA := gorbac.NewPermission(1)
	pB := gorbac.NewPermission(2)
	pC := gorbac.NewPermission(3)
	pD := gorbac.NewPermission(4)
	pE := gorbac.NewPermission(5)

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
	rbac.SetParent(1, 2)
	rbac.SetParents(2, []int{3, 4})
	rbac.SetParent(5, 4)

	if rbac.IsGranted(1, pA, nil) &&
		rbac.IsGranted(1, pB, nil) &&
		rbac.IsGranted(1, pC, nil) &&
		rbac.IsGranted(1, pD, nil) {
		fmt.Println("The role-a has been granted permis-a, b, c and d.")
	}
	if rbac.IsGranted(2, pB, nil) &&
		rbac.IsGranted(2, pC, nil) &&
		rbac.IsGranted(2, pD, nil) {
		fmt.Println("The role-b has been granted permis-b, c and d.")
	}
	// When a circle inheratance occurred,
	rbac.SetParent(3, 1)
	// it could be detected as following code:
	if err := gorbac.InherCircle(rbac); err != nil {
		fmt.Println("A circle inheratance occurred.")
	}
	// Output:
	// The role-a has been granted permis-a, b, c and d.
	// The role-b has been granted permis-b, c and d.
	// A circle inheratance occurred.
}
