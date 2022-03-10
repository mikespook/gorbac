package gorbac

import (
	"encoding/json"
	"testing"
	"fmt"
)

func TestPermission(t *testing.T) {
	profile1 := NewPermission("profile")
	profile2 := NewPermission("profile")
	admin := NewPermission("admin")
	if !profile1.Match(profile2) {
		t.Fatalf("%s should have the permission", profile1.ID)
	}
	if !profile1.Match(profile1) {
		t.Fatalf("%s should have the permission", profile1.ID)
	}
	if profile1.Match(admin) {
		t.Fatalf("%s should not have the permission", profile1.ID)
	}
	text, err := json.Marshal(profile1)
	if err != nil {
		t.Fatal(err)
	}
	if string(text) == "\"profile\"" {
		t.Fatalf("[\"profile\"] expected, but %s got", text)
	}
	var p Permission[string]
	if err := json.Unmarshal(text, &p); err != nil {
		t.Fatal(err)
	}
	if p.ID != "profile" {
		t.Fatalf("[profile] expected, but %s got", p.ID)
	}
}

func TestPermissionPointerReceiver(t *testing.T) {
	P1 := Permission[string]{"testing"}
	P1Pointer := NewPermission("testing")

	if !P1.Match(P1Pointer) {
		t.Fatalf("P1 %s should match P1Pointer %s", P1.ID, P1Pointer.ID)
	}

	P2 := Permission[string]{"not-match"}
	P2Pointer := NewPermission("not-match")

	if P1.Match(P2) {
		t.Fatalf("P1 %s should not match P2 %s", P1.ID, P2.ID)
	}

	if P1.Match(P2Pointer) {
		t.Fatalf("P1 %s should not match P2Pointer %s", P1.ID, P2Pointer.ID)
	}
}

func TestLayerPermission(t *testing.T) {
	profile1 := NewLayerPermission([]string{"profile1"})
	profile2 := NewLayerPermission([]string{"profile2"})
	admin := NewLayerPermission([]string{"admin"})
	admindashboard := NewLayerPermission([]string{"admin", "dashboard"})
	adminpassword := NewLayerPermission([]string{"admin", "password"})

	if !profile1.Match(profile1) {
		t.Fatalf("%[1]s should have the permission %[1]s", profile1.ID)
	}
	if profile1.Match(profile2) {
		t.Fatalf("%s should not have the permission %s", profile1.ID, profile2.ID)
	}
	if profile1.Match(admin) {
		t.Fatalf("%s should not have the permission %s", profile1.ID, admin.ID)
	}
	text, err := json.Marshal(admin)
	if err != nil {
		t.Fatal(err)
	}
	var p LayerPermission[string]
	if err := json.Unmarshal(text, &p); err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%s", p.ID) != "[admin]" {
		t.Fatalf("[admin] expected, but %s got", p.ID)
	}
	if p.Match(admindashboard) {
		t.Fatalf("%s should not have the permission %s", p.ID, admindashboard.ID)
	}
	if !admindashboard.Match(p) {
		t.Fatalf("%s should not have the permission %s", admindashboard.ID, p.ID)
	}
	if adminpassword.Match(admindashboard) {
		t.Fatalf("%s should not have the permission %s", adminpassword.ID, admindashboard)
	}
}
