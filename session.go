package gorbac

type SessionFactoryFunc func(string) Session

type Session interface {
	SetRbac(*Rbac)
	AddRole(string) bool
	IsGranted(string, AssertionFunc) bool
}

func NewSessionWithFactory(name string, rbac *Rbac, factory SessionFactoryFunc) Session {
	session := factory(name)
	session.SetRbac(rbac)
	return session
}

type BaseSession struct {
	Role
	rbac *Rbac
}

func NewSession(name string, rbac *Rbac) Session {
	session := &BaseSession{
		rbac:  rbac,
		Role: NewBaseRole(name),
	}
	return session
}

func (session *BaseSession) SetRbac(rbac *Rbac) {
	session.rbac = rbac
}

func (session *BaseSession) AddRole(name string) bool {
	role := session.rbac.GetRole(name)
	if role == nil {
		return false
	}
	session.Role.AddChild(role)
	return true
}

func (session *BaseSession) IsGranted(permission string, assert AssertionFunc) bool {
	if assert != nil && !assert(session.GetName(), permission, session.rbac) {
		return false
	}
	return session.HasPermission(permission)
}
