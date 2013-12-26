package gorbac

import (
	"testing"
)

const (
	FOO = "foo"
	BAR = "bar"
	OK = "ok"
	NOTEXISTS = "not-exists"
)

func TestRole(t *testing.T) {
	role := NewRole(FOO)
	if role.getName() != FOO {
		t.Errorf("`%s` expected, `%s` got.", FOO, role.getName())
	}
	if role.hasPermission(NOTEXISTS) {
		t.Errorf("`%s` should not have permission `%s`.", role.getName(), NOTEXISTS)
	}
	role.addPermission(OK)
	if !role.hasPermission(OK) {
		t.Errorf("`%s` should have permission `%s`.", role.getName(), OK)
	}
	child := NewRole(BAR)
	role.addChild(child)
	if child.getParent() != role {
		t.Errorf("`%s`'s parent should be `%s`.", child.getName(), OK)
	}
}
