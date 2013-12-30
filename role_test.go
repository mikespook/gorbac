package gorbac

import (
	"testing"
)

const (
	A         = "a"
	B         = "b"
	C         = "c"
	OK        = "ok"
	NOTEXISTS = "not-exists"
	INHERIT   = "inherit"
)

func TestRole(t *testing.T) {
	a := NewRole(A)
	if a.GetName() != A {
		t.Errorf("`%s` expected, `%s` got.", A, a.GetName())
	}
	if a.HasPermission(NOTEXISTS) {
		t.Errorf("`%s` should not have permission `%s`.", a.GetName(), NOTEXISTS)
	}
	a.AddPermission(OK)
	if !a.HasPermission(OK) {
		t.Errorf("`%s` should have permission `%s`.", a.GetName(), OK)
	}
	b := NewRole(B)
	a.AddChild(b)
	if b.GetParent() != a {
		t.Errorf("`%s`'s parent should be `%s`.", b.GetName(), OK)
	}
	c := NewRole(C)
	c.AddPermission(INHERIT)
	b.AddChild(c)
	if !a.HasPermission(INHERIT) {
		t.Errorf("`%s` should have permission `%s`.", a.GetName(), INHERIT)
	}
}

func BenchmarkHasPermission(b *testing.B) {
	ra := NewRole(A)
	rb := NewRole(B)
	ra.AddChild(rb)
	rc := NewRole(C)
	rc.AddPermission(INHERIT)
	rb.AddChild(rc)
	for i := 0; i < b.N; i++ {
		if !ra.HasPermission(INHERIT) {
			b.Errorf("`%s` should have permission `%s`.", ra.GetName(), INHERIT)
		}
	}
}
