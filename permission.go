package gorbac

import (
	"strings"
)

// Permission exports `Id` and `Match`
type Permission interface {
	ID() string
	Match(Permission) bool
}

// Permissions is a map
type Permissions map[string]Permission

// StdPermission only checks if the Ids are fully matching.
type StdPermission struct {
	IDStr string
}

// NewStdPermission returns a Permission instance with `id`
func NewStdPermission(id string) Permission {
	return &StdPermission{id}
}

// ID returns the identity of permission
func (p *StdPermission) ID() string {
	return p.IDStr
}

// Match another permission
func (p *StdPermission) Match(a Permission) bool {
	return p.IDStr == a.ID()
}

// LayerPermission firstly checks the Id of permission.
// If the Id is matched, it can be consIdered having the permission.
// Otherwise, it checks every layers of permission.
// A role which has an upper layer granted, will be granted sub-layers permissions.
type LayerPermission struct {
	IDStr string `json:"id"`
	Sep   string `json:"sep"`
}

// NewLayerPermission returns an instance of layered permission with `id`
func NewLayerPermission(id string) Permission {
	return &LayerPermission{id, ":"}
}

// ID returns the identity of permission
func (p *LayerPermission) ID() string {
	return p.IDStr
}

// Match another permission
func (p *LayerPermission) Match(a Permission) bool {
	if p.IDStr == a.ID() {
		return true
	}
	q, ok := a.(*LayerPermission)
	if !ok {
		return false
	}
	players := strings.Split(p.IDStr, p.Sep)
	qlayers := strings.Split(q.IDStr, q.Sep)
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
