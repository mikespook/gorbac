package gorbac

type Role interface {
	GetName() string
	AddPermission(string)
	HasPermission(string) bool
	RevokePermission(string)
	AddParent(string)
	RemoveParent(string)
	Reset()
}

func newBaseRole(rbac *Rbac, name string) Role {
	role := &baseRole{
		rbac: rbac,
		name:        name,
		permissions: make(map[string]bool, bufferSize),
		parents:    make(map[string]bool, bufferSize),
	}
	return role
}

type baseRole struct {
	rbac	*Rbac
	name        string
	permissions map[string]bool
	parents    map[string]bool
}

func (role *baseRole) GetName() string {
	return role.name
}

func (role *baseRole) AddPermission(permission string) {
	role.permissions[permission] = true
}

func (role *baseRole) HasPermission(permission string) bool {
	if permit, ok := role.permissions[permission]; ok {
		return permit
	}
	for pname, _ := range role.parents {
		if parent := role.rbac.GetRole(pname); parent != nil {
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
