goRBAC 
======

[![Build Status](https://travis-ci.org/mikespook/gorbac.png?branch=master)](https://travis-ci.org/mikespook/gorbac)
[![GoDoc](https://godoc.org/github.com/mikespook/gorbac?status.png)](https://godoc.org/github.com/mikespook/gorbac)
[![Coverage Status](https://coveralls.io/repos/github/mikespook/gorbac/badge.svg?branch=master)](https://coveralls.io/github/mikespook/gorbac?branch=master)

goRBAC provides a lightweight role-based access control implementation
in Golang.

For the purposes of this package:

	* an identity has one or more roles.
	* a role requests access to a permission.
	* a permission is given to a role.

Thus, RBAC has the following model:

	* many to many relationship between identities and roles.
	* many to many relationship between roles and permissions.
	* roles can have parent roles (inheriting permissions).

Version
=======

The current version of goRBAC uses Go generics (Go 1.18+) and is under active development.

Install
=======

Install the package:

```bash
$ go get github.com/mikespook/gorbac/v3
```

Usage
=====

Although you can adjust the RBAC instance anytime and it's absolutely safe, the library is designed for use with two phases:

1. Preparing

2. Checking

Preparing
---------

Import the library:

```go
import "github.com/mikespook/gorbac/v3"
```

Get a new instance of RBAC (using string as the ID type):

```go
rbac := gorbac.New[string]()
```

Get some new roles:

```go
rA := gorbac.NewRole("role-a")
rB := gorbac.NewRole("role-b")
rC := gorbac.NewRole("role-c")
rD := gorbac.NewRole("role-d")
rE := gorbac.NewRole("role-e")
```

Get some new permissions:

```go
pA := gorbac.NewPermission("permission-a")
pB := gorbac.NewPermission("permission-b")
pC := gorbac.NewPermission("permission-c")
pD := gorbac.NewPermission("permission-d")
pE := gorbac.NewPermission("permission-e")
```

Add the permissions to roles:

```go
rA.Assign(pA)
rB.Assign(pB)
rC.Assign(pC)
rD.Assign(pD)
rE.Assign(pE)
```

Also, you can implement `gorbac.Permission` for your own data structure.

After initialization, add the roles to the RBAC instance:

```go
rbac.Add(rA)
rbac.Add(rB)
rbac.Add(rC)
rbac.Add(rD)
rbac.Add(rE)
```

And set the inheritance:

```go
rbac.SetParent("role-a", "role-b")
rbac.SetParents("role-b", []string{"role-c", "role-d"})
rbac.SetParent("role-e", "role-d")
```

Checking
--------

Checking the permission is easy:

```go
if rbac.IsGranted("role-a", pA, nil) &&
	rbac.IsGranted("role-a", pB, nil) &&
	rbac.IsGranted("role-a", pC, nil) &&
	rbac.IsGranted("role-a", pD, nil) {
	fmt.Println("The role-a has been granted permis-a, b, c and d.")
}
```

Advanced Checking with Assertion Functions
------------------------------------------

You can also use assertion functions for more fine-grained permission controls:

```go
assertion := func(rbac *gorbac.RBAC[string], id string, p gorbac.Permission[string]) bool {
	// Custom logic to determine if permission should be granted
	return true // or false based on your logic
}

if rbac.IsGranted("role-a", pA, assertion) {
	fmt.Println("The role-a has been granted permission-a based on the assertion.")
}
```

Utility Functions
-----------------

goRBAC provides several built-in utility functions:

### InherCircle
Detects circular inheritance in the role hierarchy:

```go
rbac.SetParent("role-c", "role-a")
if err := gorbac.InherCircle(rbac); err != nil {
	fmt.Println("A circle inheritance occurred.")
}
```

### AnyGranted
Checks if any of the specified roles have a permission:

```go
roles := []string{"role-a", "role-b", "role-c"}
if gorbac.AnyGranted(rbac, roles, pA, nil) {
	fmt.Println("At least one role has permission-a.")
}
```

### AllGranted
Checks if all of the specified roles have a permission:

```go
roles := []string{"role-a", "role-b", "role-c"}
if gorbac.AllGranted(rbac, roles, pA, nil) {
	fmt.Println("All roles have permission-a.")
}
```

### Walk
Iterates through all roles in the RBAC instance:

```go
handler := func(r gorbac.Role[string], parents []string) error {
	fmt.Printf("Role: %s, Parents: %v\n", r.ID, parents)
	return nil
}
gorbac.Walk(rbac, handler)
```

Custom Types
------------

goRBAC supports custom types for role and permission IDs through Go generics:

```go
// Using integer IDs
rbacInt := gorbac.New[int]()
role1 := gorbac.NewRole(1)
permission1 := gorbac.NewPermission(100)

// Using custom struct IDs
type RoleID struct {
	Name string
	Type string
}

rbacStruct := gorbac.New[RoleID]()
roleCustom := gorbac.NewRole(RoleID{Name: "admin", Type: "system"})
permissionCustom := gorbac.NewPermission(RoleID{Name: "read", Type: "data"})
```

Persistence
-----------

The most asked question is how to persist the goRBAC instance. Please check the post [HOW TO PERSIST GORBAC INSTANCE](https://mikespook.com/2017/04/how-to-persist-gorbac-instance/) for the details.


Authors
=======

 * Xing Xing <mikespook@gmail.com> [Blog](http://mikespook.com) 
[@Twitter](http://twitter.com/mikespook)

Open Source - MIT Software License
==================================

See LICENSE.
