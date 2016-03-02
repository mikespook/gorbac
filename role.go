package gorbac

import (
	"sync"
)

// Role is an interface.
// You should implement this interface for your own role structures.
type Role interface {
	Id() string
	Permit(Permission) bool
}

type Roles map[string]Role
type NewRoleFunc func(string) Role

// NewStdRole is the default role factory function.
// It matches the declaration to RoleFactoryFunc.
func NewStdRole(id string) *StdRole {
	role := &StdRole{
		IdStr:       id,
		permissions: make(Permissions),
	}
	return role
}

// StdRole is the default role implement.
// You can combine this struct into your own Role implement.
type StdRole struct {
	sync.RWMutex
	IdStr       string `json:"id"`
	permissions Permissions
}

// Name returns the role's identity name.
func (role *StdRole) Id() string {
	return role.IdStr
}

// Assign a permission to the role.
func (role *StdRole) Assign(p Permission) error {
	role.Lock()
	defer role.Unlock()
	role.permissions[p.Id()] = p
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
	delete(role.permissions, p.Id())
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
