package gorbac

const (
	ParentKey      = "parents"
	PermissionKey  = "permissions"
	NameKey        = "name"
	RankKey        = "rank"
	DescriptionKey = "description"
)

// Sometimes, a custom role structure is needed by projects.
// You should define your own role factory function for this purpuse.
type RoleFactoryFunc func(*Rbac, string) Role

// An exportable data structure
type RoleMap map[string]interface{}

// An interface can't export directly. But you can convert it into a map.
func RoleToMap(role Role) RoleMap {
	roleMap := make(RoleMap)
	roleMap[PermissionKey] = role.Permissions()
	roleMap[ParentKey] = role.Parents()
	roleMap[NameKey] = role.Name()
	roleMap[RankKey] = role.Rank()
	roleMap[DescriptionKey] = role.Description()
	return roleMap
}

// Implement this interface for your own role structure.
type Role interface {
	Rank() int
	Name() string
	Description() string
	AddRank(int)
	AddDescription(string)
	AddPermission(string)
	HasPermission(string) bool
	RevokePermission(string)
	Permissions() []string
	AddParent(string)
	RemoveParent(string)
	Parents() []string
	Reset()
}

func newBaseRole(rbac *Rbac, name string) Role {
	role := &baseRole{
		rbac:        rbac,
		name:        name,
		permissions: make(map[string]bool, bufferSize),
		parents:     make(map[string]bool, bufferSize),
	}
	return role
}

type baseRole struct {
	rbac        *Rbac
	rank        int
	name        string
	description string
	permissions map[string]bool
	parents     map[string]bool
}

func (role *baseRole) AddRank(rank int) {
	role.rank = rank
}

func (role *baseRole) AddDescription(description string) {
	role.description = description
}

func (role *baseRole) AddPermission(permission string) {
	role.permissions[permission] = true
}

func (role *baseRole) HasPermission(permission string) bool {
	if permit, ok := role.permissions[permission]; ok {
		return permit
	}
	for pname, _ := range role.parents {
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

func (role *baseRole) RevokePermission(permission string) {
	delete(role.permissions, permission)
}

func (role *baseRole) AddParent(name string) {
	role.parents[name] = true
}

func (role *baseRole) RemoveParent(name string) {
	delete(role.parents, name)
}

func (role *baseRole) Reset() {
	role.permissions = make(map[string]bool, bufferSize)
	role.parents = make(map[string]bool, bufferSize)
}

func (role *baseRole) Permissions() []string {
	result := make([]string, 0, len(role.permissions))
	for name, _ := range role.permissions {
		result = append(result, name)
	}
	return result
}

func (role *baseRole) Parents() []string {
	result := make([]string, 0, len(role.parents))
	for name, _ := range role.parents {
		result = append(result, name)
	}
	return result
}

func (role *baseRole) Rank() int {
	return role.rank
}

func (role *baseRole) Name() string {
	return role.name
}

func (role *baseRole) Description() string {
	return role.description
}
