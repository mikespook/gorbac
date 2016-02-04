// User-defined gorbac example
package main

import (
	"fmt"

	"github.com/mikespook/gorbac"
)

type myRole struct {
	*gorbac.StdRole
	Label       string
	Description string
}

func loadByName(name string) (label, description string) {
	// loading data from storages or somewhere
	return name + " for testing", "This is the description for " + name
}

func newMyRole(name string) gorbac.Role {
	// loading extra properties by `name`.
	label, desc := loadByName(name)
	role := &myRole{
		Label:       label,
		Description: desc,
	}
	role.StdRole = gorbac.NewStdRole(name).(*gorbac.StdRole)
	return role
}

func main() {
	rbac := gorbac.New()
	r1 := newMyRole("role-1")
	r2 := newMyRole("role-2")
	r3 := newMyRole("role-3")
	r4 := newMyRole("role-4")
	if err := rbac.Add(r1); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if err := rbac.Add(r2); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if err := rbac.Add(r3); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if err := rbac.Add(r4); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	if err := rbac.SetParents("role-1", []string{"role-2", "role-3"}); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	if err := rbac.SetParent("role-3", "role-4"); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	role, parents, err := rbac.Get("role-1")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if myRole, ok := role.(*myRole); ok {
		fmt.Printf("Name:\t%s\nLabel:\t%s\nDesc:\t%s\nParents:\t%s\n",
			myRole.Id(), myRole.Label, myRole.Description, parents)
	}
}
