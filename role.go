package gorbac

type Role interface {
	GetName() string
	AddPermission(string)
	HasPermission(string) bool
	AddChild(Role)
	SetParent(Role)
	GetParent() Role
}

func NewRole(name string) Role {
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
	parent      Role
}

func (role *BaseRole) GetName() string {
	return role.name
}

func (role *BaseRole) AddPermission(permission string) {
	role.permissions[permission] = true
}

func (role *BaseRole) HasPermission(permission string) bool {
	if _, ok := role.permissions[permission]; ok {
		return ok
	}
	for _, child := range role.children {
		if child.HasPermission(permission) {
			role.AddPermission(permission)
			return true
		}
	}
	return false
}

func (role *BaseRole) AddChild(child Role) {
	child.SetParent(role)
	role.children = append(role.children, child)
}

func (role *BaseRole) SetParent(parent Role) {
	role.parent = parent
}

func (role *BaseRole) GetParent() Role {
	return role.parent
}
