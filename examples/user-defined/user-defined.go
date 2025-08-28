// User-defined gorbac example
package main

import (
	"fmt"

	"github.com/mikespook/gorbac/v3"
)

// myRole is a custom role that embeds the standard gorbac.Role
// and adds additional fields
type myRole struct {
	gorbac.Role[string]  // Embed the standard role
	Label       string
	Description string
}

// NewMyRole creates a new custom role with additional properties
func NewMyRole(name string) *myRole {
	// loading extra properties by `name`.
	label, desc := loadByName(name)
	return &myRole{
		Role:        gorbac.NewRole(name),  // Create the standard role
		Label:       label,
		Description: desc,
	}
}

func loadByName(name string) (label, description string) {
	// loading data from storages or somewhere
	return name + " for testing", "This is the description for " + name
}

func main() {
	rbac := gorbac.New[string]()
	r1 := NewMyRole("role-1")
	r2 := NewMyRole("role-2")
	r3 := NewMyRole("role-3")
	r4 := NewMyRole("role-4")
	
	// Add roles to RBAC - we need to pass the embedded Role part
	if err := rbac.Add(r1.Role); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if err := rbac.Add(r2.Role); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if err := rbac.Add(r3.Role); err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	if err := rbac.Add(r4.Role); err != nil {
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
	
	// Note: In this simple example, we're not demonstrating access to the custom fields
	// In a real application, you would maintain a separate map of custom roles
	fmt.Printf("Role ID: %s\nParents: %v\n", role.ID, parents)
}

