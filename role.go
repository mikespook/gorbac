package gorbac

import (
	"sync"
)

// Roles is a map
type Roles[T comparable] map[T]Role[T]

// NewStdRole is the default role factory function.
// It matches the declaration to RoleFactoryFunc.
func NewRole[T comparable](id T) Role[T] {
	return Role[T]{
		ID:          id,
		permissions: make(Permissions[T]),
	}
}

// StdRole is the default role implement.
// You can combine this struct into your own Role implement.
// T is the type of ID
type Role[T comparable] struct {
	sync.RWMutex
	// ID is the serialisable identity of role
	ID          T `json:"id"`
	permissions Permissions[T]
}

// Assign a permission to the role.
func (role *Role[T]) Assign(p Permission[T]) error {
	role.Lock()
	role.permissions[p.ID()] = p
	role.Unlock()
	return nil
}

// Permit returns true if the role has specific permission.
func (role *Role[T]) Permit(p Permission[T]) (ok bool) {
	var zero Permission[T]
	if p == zero {
		return false
	}

	role.RLock()
	for _, rp := range role.permissions {
		if rp.Match(p) {
			ok = true
			break
		}
	}
	role.RUnlock()
	return
}

// Revoke the specific permission.
func (role *Role[T]) Revoke(p Permission[T]) error {
	role.Lock()
	delete(role.permissions, p.ID())
	role.Unlock()
	return nil
}

// Permissions returns all permissions into a slice.
func (role *Role[T]) Permissions() []Permission[T] {
	role.RLock()
	result := make([]Permission[T], 0, len(role.permissions))
	for _, p := range role.permissions {
		result = append(result, p)
	}
	role.RUnlock()
	return result
}
