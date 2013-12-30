package gorbac

const (
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

func (rbac *Rbac) AddRole(role string, parent ...string) bool {
	return false
}
func (rbac *Rbac) GetRole(role string) Role {
	return nil
}
func (rbac *Rbac) HasRole(role string) bool {
	return false
}
func (rbac *Rbac) IsGranted(role, permission string,
	assert AssertionFunc) bool {
	return false
}
