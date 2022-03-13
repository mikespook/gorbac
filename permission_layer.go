package gorbac

import (
	"strings"
)

// NewLayerPermission returns an instance of layered permission with `id`
func NewLayerPermission(id, sep string) LayerPermission {
	return LayerPermission{id, sep}
}

// LayerPermission uses string as a layered ID.
// Each layer splits by "/".
type LayerPermission struct {
	SID string `json:"id"`
	Sep string `json:"sep"`
}

// ID returns id
func (p LayerPermission) ID() string {
	return p.SID
}

// Match another permission
func (p LayerPermission) Match(parent Permission[string]) bool {
	if p.SID == parent.ID() {
		return true
	}
	q, ok := parent.(LayerPermission)
	if !ok {
		return false
	}
	players := strings.Split(p.SID, p.Sep)
	qlayers := strings.Split(q.SID, q.Sep)
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
