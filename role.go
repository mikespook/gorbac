package gorbac

type Role interface {
	GetName() string
	AddPermission(string)
	HasPermission(string) bool
	AddChild(Role)
}

func NewBaseRole(name string) Role {
	role := &BaseRole{
		name:        name,
		permissions: make(map[string]bool, bufferSize),
		children:    make([]Role, 0, bufferSize),
	}
	return role
}

type BaseRole struct {
	name        string
	permissions map[string]bool
	children    []Role
}

func (role *BaseRole) GetName() string {
	return role.name
}

func (role *BaseRole) AddPermission(permission string) {
	role.permissions[permission] = true
}

func (role *BaseRole) HasPermission(permission string) bool {
	if permit := role.permissions[permission]; permit {
		return permit
	}
	for _, child := range role.children {
		if child.HasPermission(permission) {
			return true
		}
	}
	return false
}

func (role *BaseRole) AddChild(child Role) {
	role.children = append(role.children, child)
}
