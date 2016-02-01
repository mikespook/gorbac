goRBAC 
======

__NOTE:__

1. The original master branch has been moved to the branch [v1.dev](https://github.com/mikespook/gorbac/tree/v1.dev) with stable release tag [v1.0](https://github.com/mikespook/gorbac/tree/v1.0).

2. Current master comes from the redesign branch and is under heavy construction. DO NOT USE!

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

_TO BE DONE_
	
Usage
=====

_TO BE DONE_

Authors
=======

 * Xing Xing <mikespook@gmail.com> [Blog](http://mikespook.com) 
[@Twitter](http://twitter.com/mikespook)

Open Source - MIT Software License
==================================

See LICENSE.
