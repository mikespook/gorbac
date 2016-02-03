package gorbac

import (
	"strings"
)

type NewPermissionFunc func(string) Permission

type Permission interface {
	Id() string
	Match(Permission) bool
}

type Permissions map[string]Permission

// StdPermission only checks if the Ids are fully matching.
type StdPermission struct {
	IdStr string
}

func NewStdPermission(id string) Permission {
	return &StdPermission{id}
}

func (p *StdPermission) Id() string {
	return p.IdStr
}

func (p *StdPermission) Match(a Permission) bool {
	return p.IdStr == a.Id()
}

// LayerPermission firstly checks the Id of permission.
// If the Id is matched, it can be consIdered having the permission.
// Otherwise, it checks every layers of permission.
// A role which has an upper layer granted, will be granted sub-layers permissions.
type LayerPermission struct {
	IdStr string `json:"id"`
	Sep   string `json:"sep"`
}

func NewLayerPermission(id string) Permission {
	return &LayerPermission{id, ":"}
}

func (p *LayerPermission) Id() string {
	return p.IdStr
}

func (p *LayerPermission) Match(a Permission) bool {
	if p.IdStr == a.Id() {
		return true
	}
	q, ok := a.(*LayerPermission)
	if !ok {
		return false
	}
	players := strings.Split(p.IdStr, p.Sep)
	qlayers := strings.Split(q.IdStr, q.Sep)
	// layer counts of q should be less than that of p
	if len(players) > len(qlayers) {
		return false
	}
	for k, pv := range players {
		if pv != qlayers[k] {
			return false
		}
	}
	return true
}
