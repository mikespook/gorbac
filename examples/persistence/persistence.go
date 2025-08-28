package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mikespook/gorbac/v3"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func LoadJson(filename string, v interface{}) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}

func SaveJson(filename string, v interface{}) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(v)
}

func main() {
	// map[RoleId]PermissionIds
	var jsonRoles map[string][]string
	// map[RoleId]ParentIds
	var jsonInher map[string][]string
	// Load roles information
	if err := LoadJson("roles.json", &jsonRoles); err != nil {
		log.Fatal(err)
	}
	// Load inheritance information
	if err := LoadJson("inher.json", &jsonInher); err != nil {
		log.Fatal(err)
	}
	rbac := gorbac.New[string]()
	permissions := make(map[string]gorbac.Permission[string])

	// Build roles and add them to goRBAC instance
	for rid, pids := range jsonRoles {
		role := gorbac.NewRole(rid)
		for _, pid := range pids {
			_, ok := permissions[pid]
			if !ok {
				permissions[pid] = gorbac.NewPermission(pid)
			}
			role.Assign(permissions[pid])
		}
		rbac.Add(role)
	}
	// Assign the inheritance relationship
	for rid, parents := range jsonInher {
		if err := rbac.SetParents(rid, parents); err != nil {
			log.Fatal(err)
		}
	}
	// Check if `editor` can add text
	if rbac.IsGranted("editor", permissions["add-text"], nil) {
		log.Println("Editor can add text")
	}
	// Check if `chief-editor` can add text
	if rbac.IsGranted("chief-editor", permissions["add-text"], nil) {
		log.Println("Chief editor can add text")
	}
	// Check if `photographer` can add text
	if !rbac.IsGranted("photographer", permissions["add-text"], nil) {
		log.Println("Photographer can't add text")
	}
	// Check if `nobody` can add text
	// `nobody` is not exist in goRBAC at the moment
	if !rbac.IsGranted("nobody", permissions["read-text"], nil) {
		log.Println("Nobody can't read text")
	}
	// Add `nobody` and assign `read-text` permission
	nobody := gorbac.NewRole("nobody")
	permissions["read-text"] = gorbac.NewPermission("read-text")
	nobody.Assign(permissions["read-text"])
	rbac.Add(nobody)
	// Check if `nobody` can read text again
	if rbac.IsGranted("nobody", permissions["read-text"], nil) {
		log.Println("Nobody can read text")
	}

	// Persist the change
	// map[RoleId]PermissionIds
	jsonOutputRoles := make(map[string][]string)
	// map[RoleId]ParentIds
	jsonOutputInher := make(map[string][]string)
	SaveJsonHandler := func(r gorbac.Role[string], parents []string) error {
		// WARNING: Don't use gorbac.RBAC instance in the handler,
		// otherwise it causes deadlock.
		permissions := make([]string, 0)
		for _, p := range r.Permissions() {
			permissions = append(permissions, p.ID())
		}
		jsonOutputRoles[r.ID] = permissions
		jsonOutputInher[r.ID] = parents
		return nil
	}
	if err := gorbac.Walk(rbac, SaveJsonHandler); err != nil {
		log.Fatalln(err)
	}

	// Save roles information
	if err := SaveJson("new-roles.json", &jsonOutputRoles); err != nil {
		log.Fatal(err)
	}
	// Save inheritance information
	if err := SaveJson("new-inher.json", &jsonOutputInher); err != nil {
		log.Fatal(err)
	}
}
