package gorbac_test

import (
	"fmt"
	"github.com/mikespook/gorbac"
)

func ExampleRbac() {
	// Get a new access control list
	rbac := gorbac.New()

	rc := rbac.AddRole("c", "a", "b")
	rc.AddPermission("c")
	ra := rbac.AddRole("a")
	ra.AddPermission("a")

	if rbac.IsGranted("c", "c", nil) {
		fmt.Println("c granted c")
	}
	if rbac.IsGranted("a", "c", nil) {
		fmt.Println("a granted c")
	}
	if rbac.IsGranted("b", "c", nil) {
		fmt.Println("b granted c")
	}
	if rbac.IsGranted("a", "a", nil) {
		fmt.Println("a granted a")
	}
	if rbac.IsGranted("b", "a", nil) {
		fmt.Println("b granted a")
	}
	if rbac.IsGranted("c", "a", nil) {
		fmt.Println("c granted a")
	}
	rbac.Reset()
	if rbac.IsGranted("a", "a", nil) {
		fmt.Println("a granted a")
	}
	if rbac.IsGranted("b", "a", nil) {
		fmt.Println("b granted a")
	}
	if rbac.IsGranted("c", "a", nil) {
		fmt.Println("c granted a")
	}
	// Output: c granted c
	// a granted c
	// b granted c
	// a granted a
}
