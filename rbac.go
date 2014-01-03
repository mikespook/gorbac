/*
goRBAC provides a lightweight role-based access control implementation
in Golang.

For the purposes of this package:

	* an identity has one or more roles.
	* a role requests access to a permission.
	* a permission is given to a role.

Thus, RBAC has the following model:

	* many to many relationship between identities and roles.
	* many to many relationship between roles and permissions.
	* roles can have parent roles.
*/
package gorbac

import (
	"sync"
)

const (
	bufferSize = 16
)

type AssertionFunc func(string, string, *Rbac) bool
type RoleFactoryFunc func(*Rbac, string) Role

type Rbac struct {
	sync.RWMutex
	roles   map[string]Role
	factory RoleFactoryFunc
}

func NewWithFactory(factory RoleFactoryFunc) *Rbac {
	rbac := &Rbac{
		roles:   make(map[string]Role, bufferSize),
		factory: factory,
	}
	return rbac
}

func New() *Rbac {
	return NewWithFactory(newBaseRole)
}

func (rbac *Rbac) Set(name string, permissions []string, parents []string) {
	rbac.Lock()
	defer rbac.Unlock()
	role := rbac.getRole(name)
	role.Reset()
	for _, p := range permissions {
		role.AddPermission(p)
	}
	for _, pname := range parents {
		role.AddParent(pname)
	}
	rbac.roles[name] = role
}

func (rbac *Rbac) Add(name string, permissions []string, parents []string) {
	rbac.Lock()
	defer rbac.Unlock()
	role := rbac.getRole(name)
	for _, p := range permissions {
		role.AddPermission(p)
	}
	for _, pname := range parents {
		role.AddParent(pname)
	}
	rbac.roles[name] = role
}

func (rbac *Rbac) Remove(name string) {
	rbac.Lock()
	defer rbac.Unlock()
	delete(rbac.roles, name)
}

func  (rbac *Rbac) getRole(name string) Role {
	role, ok := rbac.roles[name]
	if !ok {
		role = rbac.factory(rbac, name)
	}
	return role
}

func  (rbac *Rbac) GetRole(name string) Role {
	role, ok := rbac.roles[name]
	if !ok {
		return nil
	}
	return role
}

func (rbac *Rbac) IsGranted(name, permission string,
	assert AssertionFunc) bool {
	rbac.RLock()
	defer rbac.RUnlock()
	if assert != nil && !assert(name, permission, rbac) {
		return false
	}
	if role, ok := rbac.roles[name]; ok {
		return role.HasPermission(permission)
	}
	return false
}
