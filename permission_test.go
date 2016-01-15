package gorbac

import (
	"encoding/json"
	"testing"
)

func TestStdPermission(t *testing.T) {
	profile1 := NewStdPermission("profile")
	profile2 := NewStdPermission("profile")
	admin := NewStdPermission("admin")
	if !profile1.Match(profile2) {
		t.Errorf("%s should have the permission", profile1.Id())
	}
	if !profile1.Match(profile1) {
		t.Errorf("%s should have the permission", profile1.Id())
	}
	if profile1.Match(admin) {
		t.Errorf("%s should not have the permission", profile1.Id())
	}
	text, err := json.Marshal(profile1)
	if err != nil {
		t.Error(err)
	}
	if string(text) == "\"profile\"" {
		t.Errorf("[\"profile\"] expected, but %s got", text)
	}
	var p StdPermission
	if err := json.Unmarshal(text, &p); err != nil {
		t.Error(err)
	}
	if p.Id() != "profile" {
		t.Errorf("[profile] expected, but %s got", p.Id())
	}
}

func TestLayerPermission(t *testing.T) {
	profile1 := NewLayerPermission("profile", ":")
	profile2 := NewLayerPermission("profile", ":")
	admin := NewLayerPermission("admin", ":")
	admindashboard := NewLayerPermission("admin:dashboard", ":")
	adminpassword := NewLayerPermission("admin:password", ":")

	if !profile1.Match(profile1) {
		t.Errorf("%s should have the permission", profile1.Id())
	}
	if !profile1.Match(profile2) {
		t.Errorf("%s should have the permission", profile1.Id())
	}
	if profile1.Match(admin) {
		t.Errorf("%s should not have the permission", profile1.Id())
	}
	text, err := json.Marshal(admin)
	if err != nil {
		t.Error(err)
	}
	var p LayerPermission
	if err := json.Unmarshal(text, &p); err != nil {
		t.Error(err)
	}
	if p.Id() != "admin" {
		t.Errorf("[admin] expected, but %s got", p.Id())
	}
	if !p.Match(admindashboard) {
		t.Errorf("%s should have the permission", p.Id())
	}
	if admindashboard.Match(&p) {
		t.Errorf("%s should not have the permission", admindashboard.Id())
	}
	if adminpassword.Match(admindashboard) {
		t.Errorf("%s should not have the permission", adminpassword.Id())
	}
}
