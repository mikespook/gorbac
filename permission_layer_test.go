package gorbac

import (
	"encoding/json"
	"testing"
)

func TestLayerPermission(t *testing.T) {
	profile1 := NewLayerPermission("profile1", "::")
	profile2 := NewLayerPermission("profile2", "::")
	admin := NewLayerPermission("admin", "::")
	admindashboard := NewLayerPermission("admin::dashboard", "::")
	adminpassword := NewLayerPermission("admin::password", "::")

	if !profile1.Match(profile1) {
		t.Fatalf("`%[1]s` should have the permission `%[1]s`", profile1.ID())
	}
	if profile1.Match(profile2) {
		t.Fatalf("`%s` should not have the permission `%s`", profile1.ID(), profile2.ID())
	}
	if profile1.Match(admin) {
		t.Fatalf("`%s` should not have the permission `%s`", profile1.ID(), admin.ID())
	}
	text, err := json.Marshal(admin)
	if err != nil {
		t.Fatal(err)
	}
	var p LayerPermission
	if err := json.Unmarshal(text, &p); err != nil {
		t.Fatal(err)
	}
	if p.ID() != "admin" {
		t.Fatalf("`admin` expected, but `%s` got", p.ID())
	}
	if !p.Match(admindashboard) {
		t.Fatalf("`%s` should have the permission `%s`", p.ID(), admindashboard.ID())
	}
	if admindashboard.Match(p) {
		t.Fatalf("`%s` should not have the permission `%s`", admindashboard.ID(), p.ID())
	}
	if adminpassword.Match(admindashboard) {
		t.Fatalf("`%s` should not have the permission `%s`", adminpassword.ID(), admindashboard)
	}
}
