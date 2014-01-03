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
If the role is not existing, new will be created:
	
	rbac.Add("editor", []string{"edit.article"}, nil)	
	rbac.Set("master", []string{"del.article"}, []string{"editor"})

The mainly difference between `Add` and `Set` is: 

 * `Add` keeps the permissions and the parents which are already existed;
 * `Set` covers them with new.

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

Authors
=======

 * Xing Xing <mikespook@gmail.com> [Blog](http://mikespook.com) 
[@Twitter](http://twitter.com/mikespook)

Open Source - MIT Software License
==================================
Copyright (c) 2012 Xing Xing

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

