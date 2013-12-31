package gorbac

import (
	"testing"
)

const (
	// Roles
	RA = "a"
	RB = "b"
	RC = "c"
	// Permissions
	PA = "a"
	PB = "b"
	PC = "c"

	NOTEXISTS = "not-exists"
)

func TestRole(t *testing.T) {
	ra := NewBaseRole(RA)
	// GetName
	if ra.GetName() != RA {
		t.Errorf("`%s` expected, `%s` got.", RA, ra.GetName())
	}
	// No permission
	if ra.HasPermission(NOTEXISTS) {
		t.Errorf("`%s` should not have permission `%s`.", ra.GetName(), NOTEXISTS)
	}
	// Has permission
	ra.AddPermission(PA)
	if !ra.HasPermission(PA) {
		t.Errorf("`%s` should have permission `%s`.", ra.GetName(), PA)
	}
	// Roles compositions
	rb := NewBaseRole(RB)
	rc := NewBaseRole(RC)
	ra.AddChild(rb)
	rb.AddChild(rc)
	rb.AddPermission(PB)
	rc.AddPermission(PC)
	if !ra.HasPermission(PB) {
		t.Errorf("`%s` should have permission `%s`.", ra.GetName(), PB)
	}
	if !ra.HasPermission(PC) {
		t.Errorf("`%s` should have permission `%s`.", ra.GetName(), PC)
	}
}

func BenchmarkHasPermission(b *testing.B) {
	// Initialization
	ra := NewBaseRole(RA)
	rb := NewBaseRole(RB)
	rc := NewBaseRole(RC)
	// Composition
	ra.AddChild(rb)
	rb.AddChild(rc)
	// Permission
	ra.AddPermission(PA)
	rb.AddPermission(PB)
	rc.AddPermission(PC)
	for i := 0; i < b.N; i++ {
		if !ra.HasPermission(PA) {
			b.Errorf("`%s` should have permission `%s`.", ra.GetName(), PA)
		}
		if !ra.HasPermission(PB) {
			b.Errorf("`%s` should have permission `%s`.", ra.GetName(), PB)
		}
		if !ra.HasPermission(PC) {
			b.Errorf("`%s` should have permission `%s`.", ra.GetName(), PC)
		}
	}
}
