package gorbac

const (
	ParentKey     = "parents"
	PermissionKey = "permissions"
	NameKey       = "name"
)

// Sometimes, a custom role structure is needed by projects.
// You should define your own role factory function for this purpuse.
type RoleFactoryFunc func(*Rbac, string) Role

// An exportable data structure
type RoleMap map[string][]string

// An interface can't export directly. But you can convert it into a map.
func RoleToMap(role Role) RoleMap {
	roleMap := make(RoleMap)
	roleMap[PermissionKey] = role.Permissions()
	roleMap[ParentKey] = role.Parents()
	roleMap[NameKey] = []string{role.Name()}
	return roleMap
}

// Implement this interface for your own role structure.
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
}

func NewBaseRole(rbac *Rbac, name string) Role {
	role := &BaseRole{
		rbac:        rbac,
		name:        name,
		permissions: make(map[string]bool, bufferSize),
		parents:     make(map[string]bool, bufferSize),
	}
	return role
}

type BaseRole struct {
	rbac        *Rbac
	name        string
	permissions map[string]bool
	parents     map[string]bool
}

func (role *BaseRole) Name() string {
	return role.name
}

func (role *BaseRole) AddPermission(permission string) {
	role.permissions[permission] = true
}

func (role *BaseRole) HasPermission(permission string) bool {
	if permit, ok := role.permissions[permission]; ok {
		return permit
	}
	for pname := range role.parents {
		if parent := role.rbac.Get(pname); parent != nil {
			if parent.HasPermission(permission) {
				return true
			}
		} else {
			delete(role.parents, pname)
		}
	}
	return false
}

func (role *BaseRole) RevokePermission(permission string) {
	delete(role.permissions, permission)
}

func (role *BaseRole) AddParent(name string) {
	role.parents[name] = true
}

func (role *BaseRole) RemoveParent(name string) {
	delete(role.parents, name)
}

func (role *BaseRole) Reset() {
	role.permissions = make(map[string]bool, bufferSize)
	role.parents = make(map[string]bool, bufferSize)
}

func (role *BaseRole) Permissions() []string {
	result := make([]string, 0, len(role.permissions))
	for name := range role.permissions {
		result = append(result, name)
	}
	return result
}

func (role *BaseRole) Parents() []string {
	result := make([]string, 0, len(role.parents))
	for name := range role.parents {
		result = append(result, name)
	}
	return result
}
