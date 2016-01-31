// User-defined gorbac example
package main

import (
	"fmt"

	"gopkg.in/mikespook/gorbac.v1"
)

type myRole struct {
	gorbac.Role
	Label       string
	Description string
}

func loadByName(name string) (label, description string) {
	// loading data from storages or somewhere
	return name, "This is the description for " + name
}

func newMyRole(rbac *gorbac.RBAC, name string) gorbac.Role {
	// loading extra properties by `name`.
	label, desc := loadByName(name)
	role := &myRole{
		Role:        gorbac.NewBaseRole(rbac, name),
		Label:       label,
		Description: desc,
	}
	return role
}

func main() {
	rbac := gorbac.NewWithFactory(newMyRole)
	rbac.Add("role-1", []string{"a", "b", "c"}, nil)
	role := rbac.Get("role-1")
	if myRole, ok := role.(*myRole); ok {
		fmt.Printf("Name:\t%s\nLabel:\t%s\nDesc:\t%s\n",
			myRole.Name(), myRole.Label, myRole.Description)
	}
}
