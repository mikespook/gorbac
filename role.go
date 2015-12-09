package gorbac

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
type RoleMap map[string][]string

// RoleToMap converts interface Role into RoleMap.
func RoleToMap(role Role) RoleMap {
	roleMap := make(RoleMap)
	roleMap[PermissionKey] = role.Permissions()
	roleMap[ParentKey] = role.Parents()
	roleMap[NameKey] = []string{role.Name()}
	return roleMap
}

// Role is an interface.
// You should implement this interface for your own role structures.
type Role interface {
	Name() string
	AddPermission(string)
	HasPermission(string) bool
	RevokePermission(string)
	Permissions() []string
	AddParent(string)
	RemoveParent(string)
	Parents() []string
	Reset()
	MarshalText() ([]byte, error)
	UnmarshalText([]byte) error
}

// NewBaseRole is the default role factory function.
// It matches the declaration to RoleFactoryFunc.
func NewBaseRole(rbac *RBAC, name string) Role {
	role := &BaseRole{
		rbac:        rbac,
		name:        name,
		permissions: make(map[string]bool, bufferSize),
		parents:     make(map[string]bool, bufferSize),
	}
	return role
}

// BaseRole is the default role implement.
// You can combine this struct into your own Role implement.
type BaseRole struct {
	rbac        *RBAC
	name        string
	permissions map[string]bool
	parents     map[string]bool
}

// MarshalText encodes a BaseRole into JSON format.
func (role *BaseRole) MarshalText() (text []byte, err error) {
	return nil, nil
}

// UnmarshalText decodes a JSON format into BaseRole.
func (role *BaseRole) UnmarshalText(text []byte) error {
	return nil
}

// Name returns the role's identity name.
func (role *BaseRole) Name() string {
	return role.name
}

// AddPermission adds a permission to the role.
func (role *BaseRole) AddPermission(permission string) {
	role.permissions[permission] = true
}

// HasPermission returns true if the role has specific permission.
func (role *BaseRole) HasPermission(permission string) bool {
	if permit, ok := role.permissions[permission]; ok {
		return permit
	}
	for pname := range role.parents {
		if parent := role.rbac.get(pname); parent != nil {
			if parent.HasPermission(permission) {
				return true
			}
		} else {
			delete(role.parents, pname)
		}
	}
	return false
}

// RevokePermission remove the specific permission.
func (role *BaseRole) RevokePermission(permission string) {
	delete(role.permissions, permission)
}

// AddParent adds a parent to the role.
func (role *BaseRole) AddParent(name string) {
	role.parents[name] = true
}

// RemoveParent deletes the specific parent from the role.
func (role *BaseRole) RemoveParent(name string) {
	delete(role.parents, name)
}

// Reset cleans all permissions and parents.
func (role *BaseRole) Reset() {
	role.permissions = make(map[string]bool, bufferSize)
	role.parents = make(map[string]bool, bufferSize)
}

// Permissions returns all permissions into a slice.
func (role *BaseRole) Permissions() []string {
	result := make([]string, 0, len(role.permissions))
	for name := range role.permissions {
		result = append(result, name)
	}
	return result
}

// Parents returns all parents into a slice.
func (role *BaseRole) Parents() []string {
	result := make([]string, 0, len(role.parents))
	for name := range role.parents {
		result = append(result, name)
	}
	return result
}
