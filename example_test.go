package gorbac_test

import (
	"fmt"
	"github.com/mikespook/gorbac"
)

func ExampleRbac() {
	// Get a new access control list
	rbac := gorbac.New(true)
	fmt.Println(rbac)
	// Output: &{true}
}
