/*
Package gorbac provides a lightweight role-based access
control implementation in Golang.

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
	"errors"
	"sync"
)

var (
	ErrRoleNotExist = errors.New("Role does not exist")
	ErrRoleExist    = errors.New("Role has already existed")
)

// AssertionFunc supplies more fine-grained permission controls.
type AssertionFunc func(*RBAC, string, Permission) bool

// RBAC object, in most cases it should be used as a singleton.
type RBAC struct {
	mutex       sync.RWMutex
	roles       Roles
	permissions Permissions
	parents     map[string]map[string]struct{}
}

// New returns a RBAC structure.
// The default role structure will be used.
func New() *RBAC {
	return &RBAC{
		roles:       make(Roles),
		permissions: make(Permissions),
		parents:     make(map[string]map[string]struct{}),
	}
}

// Parents bind a role `id` with `parents`.
// If the role or any of parents is not existing,
// an error will be returned.
func (rbac *RBAC) Parents(id string, parents []string) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.roles[id]; !ok {
		return ErrRoleNotExist
	}
	for _, parent := range parents {
		if _, ok := rbac.roles[parent]; !ok {
			return ErrRoleNotExist
		}
	}
	if _, ok := rbac.parents[id]; !ok {
		rbac.parents[id] = make(map[string]struct{})
	}
	var empty struct{}
	for _, parent := range parents {
		rbac.parents[id][parent] = empty
	}
	return nil
}

// Parent bind a role `id` with a `parent`.
// If the role or any of the parent is not existing,
// an error will be returned.
func (rbac *RBAC) Parent(id string, parent string) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.roles[parent]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.parents[id]; !ok {
		rbac.parents[id] = make(map[string]struct{})
	}
	var empty struct{}
	rbac.parents[id][parent] = empty
	return nil
}

// Add a role `r`.
func (rbac *RBAC) Add(r Role) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.roles[r.Id()]; ok {
		return ErrRoleExist
	}
	rbac.roles[r.Id()] = r
	return nil
}

// Remove a role.
func (rbac *RBAC) Remove(id string) {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	delete(rbac.roles, id)
	for rid, parents := range rbac.parents {
		if rid == id {
			delete(rbac.parents, rid)
			continue
		}
		for parent, _ := range parents {
			if parent == id {
				delete(rbac.parents[rid], id)
				break
			}
		}
	}
}

// Get a role and a slice of parents id.
func (rbac *RBAC) Get(id string) (Role, []string, error) {
	rbac.mutex.RLock()
	defer rbac.mutex.RUnlock()
	return rbac.get(id)
}

func (rbac *RBAC) get(id string) (Role, []string, error) {
	r, ok := rbac.roles[id]
	if !ok {
		return nil, nil, ErrRoleNotExist
	}
	parents := make([]string, 0)
	for parent, _ := range rbac.parents[id] {
		parents = append(parents, parent)
	}
	return r, parents, nil
}

// IsGranted tests if the role `id` has Permission `p` with the condition `assert`.
func (rbac *RBAC) IsGranted(id string, p Permission, assert AssertionFunc) bool {
	rbac.mutex.RLock()
	defer rbac.mutex.RUnlock()
	return rbac.isGranted(id, p, assert)
}

func (rbac *RBAC) isGranted(id string, p Permission, assert AssertionFunc) bool {
	if assert != nil && !assert(rbac, id, p) {
		return false
	}
	if role, ok := rbac.roles[id]; ok {
		return role.HasPermission(p)
	}
	return false
}
