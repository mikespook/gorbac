package gorbac

const(
	bufferSize = 16
)

type AssertionFunc func(*Rbac) bool

type Rbac struct {
	createMissing bool
}

func New(createMissing bool) *Rbac {
	rbac := &Rbac{
		createMissing: createMissing,
	}
	return rbac
}

func (rbac *Rbac) addRole(role string, parent ... string) bool {
	return false
}
func (rbac *Rbac) getRole(role string) Role {
	return nil
}
func (rbac *Rbac) hasRole(role string) bool {
	return false
}
func (rbac *Rbac) isGranted(role, permission string,
	assert AssertionFunc) bool {
	return false
}
