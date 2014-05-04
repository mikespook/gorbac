package gorbac_test

import (
	"fmt"
	"github.com/mikespook/gorbac"
)

/*
Suppose:

The role-c is inheriting from role-a, role-b.
The role-d is individual.
The role-e is inheriting from role-c, role-d.
Every roles have thire own permissions.

Thus:

The role-c has been granted permis-a/b/c.
The role-e has been granted permis-a/b/c/d.
*/
func ExampleRbac() {
	testingcases := map[string]map[string][]string{
		"role-a": map[string][]string{
			"permissions": []string{"permis-a"},
			"parents":     nil,
		},
		"role-b": map[string][]string{
			"permissions": []string{"permis-b"},
			"parents":     nil,
		},
		"role-c": map[string][]string{
			"permissions": []string{"permis-c"},
			"parents":     []string{"role-a", "role-b"},
		},
		"role-d": map[string][]string{
			"permissions": []string{"permis-d"},
			"parents":     nil,
		},
		"role-e": map[string][]string{
			"permissions": nil,
			"parents":     []string{"role-c", "role-d"},
		},
	}
	rbac := gorbac.New()

	for role, testingcase := range testingcases {
		rbac.Add(role, testingcase["permissions"], testingcase["parents"])
	}

	if rbac.IsGranted("role-c", "permis-a", nil) &&
		rbac.IsGranted("role-c", "permis-b", nil) &&
		rbac.IsGranted("role-c", "permis-c", nil) {
		fmt.Println("The role-c has been granted permis-a/b/c.")
	}
	if rbac.IsGranted("role-e", "permis-a", nil) &&
		rbac.IsGranted("role-e", "permis-b", nil) &&
		rbac.IsGranted("role-e", "permis-c", nil) &&
		rbac.IsGranted("role-e", "permis-d", nil) {
		fmt.Println("The role-e has been granted permis-a/b/c/d.")
	}
	// Output: The role-c has been granted permis-a/b/c.
	// The role-e has been granted permis-a/b/c/d.
}
