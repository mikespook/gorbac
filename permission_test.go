package gorbac

import (
	"encoding/json"
	"testing"
)

func TestStdPermission(t *testing.T) {
	// permission a
	pa := &StdPermission{PA}
	// permission b
	pb := &StdPermission{PB}
	// another instance of permission a
	paa := &StdPermission{PA}
	// different instance, matched permission
	if !pa.Has(paa) {
		t.Errorf("%s should have the permission", pa.Name())
	}
	// same instance, matched permission
	if !pa.Has(pa) {
		t.Errorf("%s should have the permission", pa.Name())
	}
	// unmatched permission
	if pa.Has(pb) {
		t.Errorf("%s should not have the permission", pa.Name())
	}
	text, err := json.Marshal(pa)
	if err != nil {
		t.Error(err)
	}
	if string(text) == PA {
		t.Errorf("%s expected, but %s got", PA, text)
	}
	var p StdPermission
	if err := json.Unmarshal(text, &p); err != nil {
		t.Error(err)
	}
	if p.Name() != PA {
		t.Errorf("%s expected, but %s got", PA, p.Name())
	}
}

func TestLayerPermission(t *testing.T) {
	// permission a with `admin` layer
	pa := &LayerPermission{PA, []string{"admin"}}
	// permission b with `profile` layer
	pb := &LayerPermission{PB, []string{"profile"}}
	// another instance of permission a
	paa := &LayerPermission{PA, []string{"admin"}}
	// permission aaa with `admin.dashboard` layer.
	// If a role is granted `pa`, then it has `paaa` automatically.
	paaa := &LayerPermission{PC, []string{"admin", "dashboard"}}
	// If a role is granted `pa`, then it has `paab` automatically.
	paab := &LayerPermission{PD, []string{"admin", "password"}}
	// different instance, matched permission
	if !pa.Has(paa) {
		t.Errorf("%s should have the permission", pa.Name())
	}
	// same instance, matched permission
	if !pa.Has(pa) {
		t.Errorf("%s should have the permission", pa.Name())
	}
	// unmatched permission
	if pa.Has(pb) {
		t.Errorf("%s should not have the permission", pa.Name())
	}
	text, err := json.Marshal(pa)
	if err != nil {
		t.Error(err)
	}
	var p LayerPermission
	if err := json.Unmarshal(text, &p); err != nil {
		t.Error(err)
	}
	if p.Name() != PA {
		t.Errorf("%s expected, but %s got", PA, p.Name())
	}
	if !p.Has(paaa) {
		t.Errorf("%s should have the permission", p.Name())
	}
	if paaa.Has(paab) {
		t.Errorf("%s should not have the permission", paaa.Name())
	}
	if paaa.Has(pa) {
		t.Errorf("%s should not have the permission", paaa.Name())
	}
}
