goRBAC 
======

__NOTE:__

0. Current master comes from the redesign branch (to be the v2 once stabled).

1. The original master branch has been moved to the branch [v1.dev](https://github.com/mikespook/gorbac/tree/v1.dev) with stable release tag [v1.0](https://github.com/mikespook/gorbac/tree/v1.0);

2. Currently, the API may change without notice and documentation;

3. Use at your own __RISK__.

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
	* roles can have a parent role (inheriting permissions).

Version
=======

Currently, goRBAC has two versions:

[Version 1](https://github.com/mikespook/gorbac/tree/v1.dev) with the original design is stable.

While [the master branch](https://github.com/mikespook/gorbac) which intended to be Version 2 is under developing.


Install
=======

Install the package:

> $ go get github.com/mikespook/gorbac
	
Usage
=====

Despite you can adjust the RBAC instance anytime and it's absolutely safe, the library is designed for using with two phases:

1. Preparing

2. Checking

Preparing
---------

Import the library:

	import "github.com/mikespook/gorbac"

Get a new instance of RBAC:

	rbac := gorbac.New()

Get some new roles:

	rA := gorbac.NewStdRole("role-a")
	rB := gorbac.NewStdRole("role-b")
	rC := gorbac.NewStdRole("role-c")
	rD := gorbac.NewStdRole("role-d")
	rE := gorbac.NewStdRole("role-e")

Get some new permissions:

	pA := gorbac.NewStdPermission("permission-a")
	pB := gorbac.NewStdPermission("permission-b")
	pC := gorbac.NewStdPermission("permission-c")
	pD := gorbac.NewStdPermission("permission-d")
	pE := gorbac.NewStdPermission("permission-e")

Add the permissions to roles:

	rA.Assign(pA)
	rB.Assign(pB)
	rC.Assign(pC)
	rD.Assign(pD)
	rE.Assign(pE)

Also, you can implement `gorbac.Role` and `gorbac.Permission` for your own data structure.

After initailization, add the roles to the RBAC instance:

	rbac.Add(rA)
	rbac.Add(rB)
	rbac.Add(rC)
	rbac.Add(rD)
	rbac.Add(rE)

And set the inheritance:

	rbac.SetParent("role-a", "role-b")
	rbac.SetParents("role-b", []string{"role-c", "role-d"})
	rbac.SetParent("role-e", "role-d")

Checking
--------

Checking the permission is easy:

	if rbac.IsGranted("role-a", pA, nil) &&
		rbac.IsGranted("role-a", pB, nil) &&
		rbac.IsGranted("role-a", pC, nil) &&
		rbac.IsGranted("role-a", pD, nil) {
		fmt.Println("The role-a has been granted permis-a, b, c and d.")
	}


And there are some built-in util-functions: 
[InherCircle](https://godoc.org/github.com/mikespook/gorbac#InherCircle),
[AnyGranted](https://godoc.org/github.com/mikespook/gorbac#AnyGranted), 
[AllGranted](https://godoc.org/github.com/mikespook/gorbac#AllGranted). 
Please [open an issue](https://github.com/mikespook/gorbac/issues/new) 
for the new built-in requriement.

E.g.:

	rbac.SetParent("role-c", "role-a")
	if err := gorbac.InherCircle(rbac); err != nil {
		fmt.Println("A circle inheratance ocurred.")
	}

Patches
=======

__2016-03-03__

    gofmt -w -r 'AssignPermission -> Assign' .
	gofmt -w -r 'RevokePermission -> Revoke' .


Authors
=======

 * Xing Xing <mikespook@gmail.com> [Blog](http://mikespook.com) 
[@Twitter](http://twitter.com/mikespook)

Open Source - MIT Software License
==================================

See LICENSE.
