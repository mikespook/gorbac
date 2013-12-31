package gorbac

import (
	"testing"
)

func TestRbac(t *testing.T) {
	rbac := New()
	rc := rbac.AddRole(RC, RA, RB)

	ra := rbac.GetRole(RA)
	if ra == nil {
		t.Error("`Role` expected, `nil` got.")
	}
	if ra.GetName() != RA {
		t.Errorf("`%s` expected, `%s` got.", RA, ra.GetName())
	}

	rb := rbac.GetRole(RB)
	if rb == nil {
		t.Error("`Role` expected, `nil` got.")
	}
	if rb.GetName() != RB {
		t.Errorf("`%s` expected, `%s` got.", RB, rb.GetName())
	}

	rc.AddPermission(PC)

	if !rbac.IsGranted(RC, PC, nil) {
		t.Errorf("`%s` granted `%s`, `false` got.", RC, PC)
	}

	if !rbac.IsGranted(RA, PC, nil) {
		t.Errorf("`%s` granted `%s`, `false` got.", RA, PC)
	}

	if !rbac.IsGranted(RB, PC, func(name, permission string, rbac *Rbac) bool {
			return name == RB && permission == PC
		}) {
		t.Error("assert return `true`, `false` got.")
	}

	if rbac.IsGranted(RB, PC, func(name, permission string, rbac *Rbac) bool {
			return name != RB || permission != PC
		}) {
		t.Error("assert return `false`, `true` got.")
	}
}
