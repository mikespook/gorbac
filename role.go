package gorbac

import (
	"sync"
)

// Role is an interface.
// You should implement this interface for your own role structures.
type Role interface {
	ID() string
	Permit(Permission) bool
}

// Roles is a map
type Roles map[string]Role

// NewStdRole is the default role factory function.
// It matches the declaration to RoleFactoryFunc.
func NewStdRole(id string) *StdRole {
	role := &StdRole{
		IDStr:       id,
		permissions: make(Permissions),
	}
	return role
}

// StdRole is the default role implement.
// You can combine this struct into your own Role implement.
type StdRole struct {
	sync.RWMutex
	// IDStr is the identity of role
	IDStr       string `json:"id"`
	permissions Permissions
}

// ID returns the role's identity name.
func (role *StdRole) ID() string {
	return role.IDStr
}

// Assign a permission to the role.
func (role *StdRole) Assign(p Permission) error {
	role.Lock()
	defer role.Unlock()
	role.permissions[p.ID()] = p
	return nil
}

// Permit returns true if the role has specific permission.
func (role *StdRole) Permit(p Permission) bool {
	role.RLock()
	defer role.RUnlock()
	for _, rp := range role.permissions {
		if rp.Match(p) {
			return true
		}
	}
	return false
}

// Revoke the specific permission.
func (role *StdRole) Revoke(p Permission) error {
	role.Lock()
	defer role.Unlock()
	delete(role.permissions, p.ID())
	return nil
}

// Permissions returns all permissions into a slice.
func (role *StdRole) Permissions() []Permission {
	role.RLock()
	defer role.RUnlock()
	result := make([]Permission, 0, len(role.permissions))
	for _, p := range role.permissions {
		result = append(result, p)
	}
	return result
}
