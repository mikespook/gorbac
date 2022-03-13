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
	// ErrRoleNotExist occurred if a role cann't be found
	ErrRoleNotExist = errors.New("Role does not exist")
	// ErrRoleExist occurred if a role shouldn't be found
	ErrRoleExist = errors.New("Role has already existed")
	empty        = struct{}{}
)

// AssertionFunc supplies more fine-grained permission controls.
type AssertionFunc[T comparable] func(*RBAC[T], T, Permission[T]) bool

// RBAC object, in most cases it should be used as a singleton.
type RBAC[T comparable] struct {
	mutex   sync.RWMutex
	roles   Roles[T]
	parents map[T]map[T]struct{}
}

// New returns a RBAC structure.
// The default role structure will be used.
func New[T comparable]() *RBAC[T] {
	return &RBAC[T]{
		roles:   make(Roles[T]),
		parents: make(map[T]map[T]struct{}),
	}
}

// SetParents bind `parents` to the role `id`.
// If the role or any of parents is not existing,
// an error will be returned.
func (rbac *RBAC[T]) SetParents(id T, parents []T) error {
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
		rbac.parents[id] = make(map[T]struct{})
	}
	for _, parent := range parents {
		rbac.parents[id][parent] = empty
	}
	return nil
}

// GetParents return `parents` of the role `id`.
// If the role is not existing, an error will be returned.
// Or the role doesn't have any parents,
// a nil slice will be returned.
func (rbac *RBAC[T]) GetParents(id T) ([]T, error) {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.roles[id]; !ok {
		return nil, ErrRoleNotExist
	}
	ids, ok := rbac.parents[id]
	if !ok {
		return nil, nil
	}
	var parents []T
	for parent := range ids {
		parents = append(parents, parent)
	}
	return parents, nil
}

// SetParent bind the `parent` to the role `id`.
// If the role or the parent is not existing,
// an error will be returned.
func (rbac *RBAC[T]) SetParent(id T, parent T) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.roles[parent]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.parents[id]; !ok {
		rbac.parents[id] = make(map[T]struct{})
	}
	var empty struct{}
	rbac.parents[id][parent] = empty
	return nil
}

// RemoveParent unbind the `parent` with the role `id`.
// If the role or the parent is not existing,
// an error will be returned.
func (rbac *RBAC[T]) RemoveParent(id T, parent T) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.roles[parent]; !ok {
		return ErrRoleNotExist
	}
	delete(rbac.parents[id], parent)
	return nil
}

// Add a role `r`.
func (rbac *RBAC[T]) Add(r Role[T]) (err error) {
	rbac.mutex.Lock()
	if _, ok := rbac.roles[r.ID]; !ok {
		rbac.roles[r.ID] = r
	} else {
		err = ErrRoleExist
	}
	rbac.mutex.Unlock()
	return
}

// Remove the role by `id`.
func (rbac *RBAC[T]) Remove(id T) (err error) {
	rbac.mutex.Lock()
	if _, ok := rbac.roles[id]; ok {
		delete(rbac.roles, id)
		for rid, parents := range rbac.parents {
			if rid == id {
				delete(rbac.parents, rid)
				continue
			}
			for parent := range parents {
				if parent == id {
					delete(rbac.parents[rid], id)
					break
				}
			}
		}
	} else {
		err = ErrRoleNotExist
	}
	rbac.mutex.Unlock()
	return
}

// Get the role by `id` and a slice of its parents id.
func (rbac *RBAC[T]) Get(id T) (r Role[T], parents []T, err error) {
	rbac.mutex.RLock()
	var ok bool
	if r, ok = rbac.roles[id]; ok {
		for parent := range rbac.parents[id] {
			parents = append(parents, parent)
		}
	} else {
		err = ErrRoleNotExist
	}
	rbac.mutex.RUnlock()
	return
}

// IsGranted tests if the role `id` has Permission `p` with the condition `assert`.
func (rbac *RBAC[T]) IsGranted(id T, p Permission[T],
	assert AssertionFunc[T]) (ok bool) {
	rbac.mutex.RLock()
	ok = rbac.isGranted(id, p, assert)
	rbac.mutex.RUnlock()
	return
}

func (rbac *RBAC[T]) isGranted(id T, p Permission[T],
	assert AssertionFunc[T]) bool {
	if assert != nil && !assert(rbac, id, p) {
		return false
	}
	return rbac.recursionCheck(id, p)
}

func (rbac *RBAC[T]) recursionCheck(id T, p Permission[T]) bool {
	if role, ok := rbac.roles[id]; ok {
		if role.Permit(p) {
			return true
		}
		if parents, ok := rbac.parents[id]; ok {
			for pID := range parents {
				if _, ok := rbac.roles[pID]; ok {
					if rbac.recursionCheck(pID, p) {
						return true
					}
				}
			}
		}
	}
	return false
}
