package gorbac

import (
	"testing"
)

var rbac *Rbac

func init() {
	rbac = New(true)
}

func TestFoobar(t *testing.T) {
	t.Log("Foobar")
}
