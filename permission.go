package gorbac

import "encoding/json"

type Permission interface {
	Name() string
	Has(Permission) bool
	MarshalText() ([]byte, error)
	UnmarshalText([]byte) error
}

// StdPermission checks if the name of permission matched.
type StdPermission struct {
	name string
}

func (p *StdPermission) Name() string {
	return p.name
}

func (p *StdPermission) Has(a Permission) bool {
	return p.name == a.Name()
}

func (p *StdPermission) MarshalText() (text []byte, err error) {
	return []byte(p.name), nil
}

func (p *StdPermission) UnmarshalText(text []byte) error {
	p.name = string(text)
	return nil
}

// LayerPermission firstly checks the name of permission.
// If the name is matched, it can be considered having the permission.
// Otherwise, it checks every layers of permission.
// A role which has an upper layer granted, will be granted sub-layers permissions.
type LayerPermission struct {
	name   string
	layers []string
}

func (p *LayerPermission) Name() string {
	return p.name
}

func (p *LayerPermission) Has(a Permission) bool {
	if p.name == a.Name() {
		return true
	}
	alp, ok := a.(*LayerPermission)
	// layers of a should be less than that of p
	if !ok {
		return false
	}
	if len(p.layers) > len(alp.layers) {
		return false
	}
	for k, pv := range p.layers {
		if pv != alp.layers[k] {
			return false
		}
	}
	return true
}

func (p *LayerPermission) MarshalText() (text []byte, err error) {
	v := map[string]interface{}{
		"name":  p.name,
		"layer": p.layers,
	}
	return json.Marshal(v)
}

func (p *LayerPermission) UnmarshalText(text []byte) error {
	var data map[string]interface{}
	if err := json.Unmarshal(text, &data); err != nil {
		return err
	}
	var ok bool
	if p.name, ok = data["name"].(string); !ok {
		return ErrUnmarshal
	}
	layer, ok := data["layer"].([]interface{})
	if !ok {
		return ErrUnmarshal
	}
	for _, v := range layer {
		l, ok := v.(string)
		if !ok {
			return ErrUnmarshal
		}
		p.layers = append(p.layers, l)
	}
	return nil
}
