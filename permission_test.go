package gorbac

import (
	"encoding/json"
	"testing"
)

func TestBasicPermission(t *testing.T) {
	profile1 := NewPermission("profile")
	profile2 := NewPermission("profile")
	admin := NewPermission("admin")
	if !profile1.Match(profile2) {
		t.Fatalf("%s should have the permission", profile1.ID())
	}
	if !profile1.Match(profile1) {
		t.Fatalf("%s should have the permission", profile1.ID())
	}
	if profile1.Match(admin) {
		t.Fatalf("%s should not have the permission", profile1.ID())
	}
	text, err := json.Marshal(profile1)
	if err != nil {
		t.Fatal(err)
	}
	if string(text) == "\"profile\"" {
		t.Fatalf("[\"profile\"] expected, but [%s] got", text)
	}
	var p StdPermission[string]
	if err := json.Unmarshal(text, &p); err != nil {
		t.Fatal(err)
	}
	if p.ID() != "profile" {
		t.Fatalf("[profile] expected, but [%s] got", p.ID())
	}
}

func TestBasicPermissionPointerReceiver(t *testing.T) {
	P1 := StdPermission[string]{"testing"}
	P1Pointer := NewPermission("testing")

	if !P1.Match(P1Pointer) {
		t.Fatalf("P1 %s should match P1Pointer %s", P1.ID(), P1Pointer.ID())
	}

	P2 := StdPermission[string]{"not-match"}
	P2Pointer := NewPermission("not-match")

	if P1.Match(P2) {
		t.Fatalf("P1 %s should not match P2 %s", P1.ID(), P2.ID())
	}

	if P1.Match(P2Pointer) {
		t.Fatalf("P1 %s should not match P2Pointer %s", P1.ID(), P2Pointer.ID())
	}
}
