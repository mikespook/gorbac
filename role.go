package gorbac

type Role interface {
	getName() string
	addPermission(string)
	hasPermission(string) bool
	addChild(Role)
	setParent(Role)
	getParent() Role
}

func NewRole(name string) Role {
	role := &innerRole{
		name: name,
		permissions: make(map[string]bool, bufferSize),
		children: make([]Role, 0, bufferSize),
	}
	return role
}

type innerRole struct {
	name string
	permissions map[string]bool
	children []Role
	parent Role
}

func (role *innerRole) getName() string {
	return role.name
}

func (role *innerRole) addPermission(permission string) {
	role.permissions[permission] = true
}

func (role *innerRole) hasPermission(permission string) bool {
	_, ok := role.permissions[permission]
	return ok
}

func (role *innerRole) addChild(child Role) {
	child.setParent(role)
	role.children = append(role.children, child)
}

func (role *innerRole) setParent(parent Role) {
	role.parent = parent
}

func (role *innerRole) getParent() Role {
	return role.parent
}
