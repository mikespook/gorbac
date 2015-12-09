package gorbac

import "encoding/json"

const (
	// ParentKey exports parents into RoleMap
	ParentKey = "parents"
	// PermissionKey exports permissions into RoleMap
	PermissionKey = "permissions"
	// NameKey exports name into RoleMap
	NameKey = "name"
)

// RoleFactoryFunc is used for a custom role structure.
// You could define your own role factory function through
// this factory function.
type RoleFactoryFunc func(*RBAC, string) Role

// RoleMap exports roles data.
type RoleMap map[string]interface{}

// RoleToMap converts interface Role into RoleMap.
func RoleToMap(role Role) RoleMap {
	roleMap := make(RoleMap)
	//	roleMap[PermissionKey] = role.permissions()
	roleMap[ParentKey] = role.Parents()
	roleMap[NameKey] = []string{role.Name()}
	return roleMap
}

// Role is an interface.
// You should implement this interface for your own role structures.
type Role interface {
	Name() string
	AddPermission(Permission)
	HasPermission(Permission) bool
	RevokePermission(Permission)
	AddParent(string)
	RemoveParent(string)
	Reset()

	Parents() []string

	MarshalText() ([]byte, error)
	UnmarshalText([]byte) error
}

type Parents map[string]interface{}

// NewStdRole is the default role factory function.
// It matches the declaration to RoleFactoryFunc.
func NewStdRole(rbac *RBAC, name string) Role {
	role := &StdRole{
		rbac:        rbac,
		name:        name,
		permissions: make(Permissions),
		parents:     make(Parents),
	}
	return role
}

// StdRole is the default role implement.
// You can combine this struct into your own Role implement.
type StdRole struct {
	rbac        *RBAC
	name        string
	permissions Permissions
	parents     Parents
}

// MarshalText encodes a StdRole into JSON format.
func (role *StdRole) MarshalText() (text []byte, err error) {
	data := map[string]interface{}{
		"name":        role.name,
		"parents":     role.Parents(),
		"permissions": role._permissions(),
	}
	return json.Marshal(data)
}

// UnmarshalText decodes a JSON format into StdRole.
func (role *StdRole) UnmarshalText(text []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(text, &data); err != nil {
		return err
	}
	var ok bool
	if role.name, ok = data["name"].(string); !ok {
		return ErrUnmarshal
	}
	parents, ok := data["parents"].([]interface{})
	if !ok {
		return ErrUnmarshal
	}
	for _, v := range parents {
		if name, ok := v.(string); ok {
			role.parents[name] = nil
		}
	}
	permissions, ok := data["permissions"].([]interface{})
	if !ok {
		return ErrUnmarshal
	}
	for _, v := range permissions {
		p, ok := v.(Permission)
		if ok {
			role.permissions[p.Name()] = p
		}
	}
	return nil
}

// Name returns the role's identity name.
func (role *StdRole) Name() string {
	return role.name
}

// AddPermission adds a permission to the role.
func (role *StdRole) AddPermission(p Permission) {
	role.permissions[p.Name()] = p
}

// HasPermission returns true if the role has specific permission.
func (role *StdRole) HasPermission(p Permission) bool {
	for _, rp := range role.permissions {
		if rp.Has(p) {
			return true
		}
	}
	for pname := range role.parents {
		if parent := role.rbac.get(pname); parent != nil {
			if parent.HasPermission(p) {
				return true
			}
		} else {
			delete(role.parents, pname)
		}
	}
	return false
}

// RevokePermission remove the specific permission.
func (role *StdRole) RevokePermission(p Permission) {
	delete(role.permissions, p.Name())
}

// AddParent adds a parent to the role.
func (role *StdRole) AddParent(name string) {
	role.parents[name] = nil
}

// RemoveParent deletes the specific parent from the role.
func (role *StdRole) RemoveParent(name string) {
	delete(role.parents, name)
}

// Reset cleans all permissions and parents.
func (role *StdRole) Reset() {
	role.permissions = make(map[string]Permission)
	role.parents = make(map[string]interface{})
}

// Permissions returns all permissions into a slice.
func (role *StdRole) _permissions() []Permission {
	result := make([]Permission, 0, len(role.permissions))
	for _, p := range role.permissions {
		result = append(result, p)
	}
	return result
}

// Parents returns all parents into a slice.
func (role *StdRole) Parents() []string {
	result := make([]string, 0, len(role.parents))
	for name := range role.parents {
		result = append(result, name)
	}
	return result
}
