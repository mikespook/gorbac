package gorbac

import (
	"testing"
)

func TestSession(t *testing.T) {
	rbac := New()
	rbac.AddRole(RB)
	rc := rbac.AddRole(RC, RA)
	rc.AddPermission(PC)

	session := NewSession("session", rbac)
	if !session.AddRole(RA) {
		t.Errorf("AddRole(%s) assert return `true`, `false` got.", RA)
	}

	if session.AddRole(NOTEXISTS) {
		t.Errorf("AddRole(%s) assert return `false`, `true` got.", NOTEXISTS)
	}
	session.AddRole(RB)

	if !session.IsGranted(PC, nil) {
		t.Errorf("IsGranted(%s) assert return `true`, `false` got.", PC)
	}

	if session.IsGranted(PA, nil) {
		t.Errorf("IsGranted(%s) assert return `false`, `true` got.", PA)
	}
}
