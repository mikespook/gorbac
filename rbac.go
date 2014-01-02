package gorbac

const (
	bufferSize = 16
)

type AssertionFunc func(string, string, *Rbac) bool
type RoleFactoryFunc func(string) Role

type Rbac struct {
	roles   map[string]Role
	factory RoleFactoryFunc
}

func NewWithFactory(factory RoleFactoryFunc) *Rbac {
	rbac := &Rbac{
		roles:   make(map[string]Role, bufferSize),
		factory: factory,
	}
	return rbac
}

func New() *Rbac {
	rbac := &Rbac{
		roles:   make(map[string]Role, bufferSize),
		factory: NewBaseRole,
	}
	return rbac
}

func (rbac *Rbac) AddRole(name string, parents ...string) (role Role) {
	role = rbac.GetRole(name)
	if role == nil {
		role = rbac.factory(name)
	}
	for _, pn := range parents {
		pr := rbac.AddRole(pn)
		pr.AddChild(role)
	}
	rbac.roles[name] = role
	return
}

func (rbac *Rbac) GetRole(name string) Role {
	return rbac.roles[name]
}

func (rbac *Rbac) IsGranted(name, permission string,
	assert AssertionFunc) bool {
	if assert != nil && !assert(name, permission, rbac) {
		return false
	}
	if role := rbac.GetRole(name); role != nil {
		return role.HasPermission(permission)
	}
	return false
}

func (rbac *Rbac) Reset() {
	rbac.roles = make(map[string]Role, bufferSize)
}
