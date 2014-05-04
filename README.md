goRBAC
======

[![Build Status](https://travis-ci.org/mikespook/gorbac.png?branch=master)](https://travis-ci.org/mikespook/gorbac)
[![GoDoc](https://godoc.org/github.com/mikespook/gorbac?status.png)](https://godoc.org/github.com/mikespook/gorbac)

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

Install
=======

Install the package:

> $ go get github.com/mikespook/gorbac
	
Usage
=====

Import the package:

	import github.com/mikespook/gorbac

Get a goRBAC instance:
	
	rbac := gorbac.New()

Specified permissions and parent roles for a role.
If the role is not existing, new one will be created:
	
	rbac.Add("editor", []string{"edit.article"}, nil)	
	rbac.Set("master", []string{"del.article"}, []string{"editor"})

The mainly difference between `Add` and `Set` is: 

 * `Add` keeps original permissions and parents which are already existed;
 * `Set` covers them with new permissions and parents.

Remove a role:

	rbac.Remove("guest")

Get a role for more fine-grained controls:

	rbac.Get("admin")

Check if a role has a permission:
	
	rbac.IsGranted("editor", "edit.article", nil)

The 3rd param, Assertion function is used for more fine-grained testing:

	rbac.IsGranted("editor", "edit.article", 
		func(role, permission string, rbac* Rbac) bool {
			return article.Owner == User.Id
	})

Revoke a permission from a role:

	rbac.Get("master").RevokePermission("del.article")

Remove a role's parent:

	rbac.Get("editor").RemoveParent("auth-user")

In a real case, it is good for checking if a role existed:

	if role := rbac.Get("not-exists"); role == nil {
		// Not exists;
	} else {
		// Exists. 	
	}

`Dump` and `Restore` help for data persistence:

	m := rbac.Dump()
	data, err := json.Marshal(m)
	// Handling error or save data

	var m gorbac.Map
	err := json.Unmarshal(data, &m)
	rbac = gorbac.Restore(m)

For more details, please see [example_test.go](https://github.com/mikespook/gorbac/blob/master/example_test.go) and [example/http](https://github.com/mikespook/gorbac/tree/master/examples/http).

Authors
=======

 * Xing Xing <mikespook@gmail.com> [Blog](http://mikespook.com) 
[@Twitter](http://twitter.com/mikespook)

Open Source - MIT Software License
==================================

See LICENSE.
