package gorbac

// Permission exports `Id` and `Match`
type Permission[K comparable] struct {
	ID K
}

// Match another permission
func (p Permission[comparable]) Match(a Permission[comparable]) bool {
	return p.ID == a.ID
}

func NewPermission[K comparable](id K) Permission[K] {
	return Permission[K]{id}
}


// Permissions is a map
type Permissions[K comparable] map[K]Permission[K]

// LayerPermission uses []K as a layered ID.
// Lower layers have higher privileges.
type LayerPermission[K comparable] struct {
	ID []K `json:"id"`
}

// NewLayerPermission returns an instance of layered permission with `id`
func NewLayerPermission[K comparable](id []K) LayerPermission[K] {
	return LayerPermission[K]{id}
}

// Match another permission
func (p *LayerPermission[comparable]) Match(parent LayerPermission[comparable]) bool {
	// The layer number of parent has to be bigger than p,
	// otherwise, parent is sub-permission of p.
	//
	// e.g. a = [1 ,2]; b = [1, 2, 3];
	// It means b is a sub permission of a.
	// i.e. b.Match(a) is true, a.Match(b) is false.
	if len(p.ID) < len(parent.ID) {
		return false
	}
	// p longer than or equal to parent
	for k, pv := range parent.ID {
		// In complex case, e.g. a = [1, 2, 3]; b = [1, 3, 4];
		// The permission is miss-matched on the second layer.
		if pv != p.ID[k] {
			return false
		}
	}
	return true
}
