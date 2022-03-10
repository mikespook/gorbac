package gorbac

import (
	"sync"
)

// Roles is a map
type Roles[K comparable] map[K]Role[K]

// NewStdRole is the default role factory function.
// It matches the declaration to RoleFactoryFunc.
func NewRole[K comparable](id K) Role[K] {
	return Role[K]{
		ID:       id,
		permissions: make(Permissions[K]),
	}
}

// StdRole is the default role implement.
// You can combine this struct into your own Role implement.
type Role[K comparable] struct {
	sync.RWMutex
	// IDStr is the identity of role
	ID       K `json:"id"`
	permissions Permissions[K]
}

// Assign a permission to the role.
func (role *Role[K]) Assign(p Permission[K]) error {
	role.Lock()
	role.permissions[p.ID] = p
	role.Unlock()
	return nil
}

// Permit returns true if the role has specific permission.
func (role *Role[K]) Permit(p Permission[K]) (ok bool) {
	var zero Permission[K]
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
func (role *Role[K]) Revoke(p Permission[K]) error {
	role.Lock()
	delete(role.permissions, p.ID)
	role.Unlock()
	return nil
}

// Permissions returns all permissions into a slice.
func (role *Role[K]) Permissions() []Permission[K] {
	role.RLock()
	result := make([]Permission[K], 0, len(role.permissions))
	for _, p := range role.permissions {
		result = append(result, p)
	}
	role.RUnlock()
	return result
}
