package gorbac

import (
	"sync"
)

// Role is an interface.
// You should implement this interface for your own role structures.
type Role interface {
	Id() string
	AddPermission(Permission) error
	HasPermission(Permission) bool
	RevokePermission(Permission) error
	Permissions() []Permission
}

type Roles map[string]Role
type NewRoleFunc func(string) Role

// NewStdRole is the default role factory function.
// It matches the declaration to RoleFactoryFunc.
func NewStdRole(id string) Role {
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

// AddPermission adds a permission to the role.
func (role *StdRole) AddPermission(p Permission) error {
	role.Lock()
	defer role.Unlock()
	role.permissions[p.Id()] = p
	return nil
}

// HasPermission returns true if the role has specific permission.
func (role *StdRole) HasPermission(p Permission) bool {
	role.RLock()
	defer role.RUnlock()
	for _, rp := range role.permissions {
		if rp.Match(p) {
			return true
		}
	}
	return false
}

// RevokePermission remove the specific permission.
func (role *StdRole) RevokePermission(p Permission) error {
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
