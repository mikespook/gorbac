package gorbac

import (
	"bytes"
	"encoding/json"
	"strings"
)

type Permission interface {
	Id() string
	Match(Permission) bool
	MarshalText() ([]byte, error)
	UnmarshalText([]byte) error
}

type Permissions map[string]Permission

func (ps Permissions) MarshalText() (text []byte, err error) {
	return json.Marshal(ps)
}

func (ps Permissions) UnmarshalText(text []byte) error {
	return json.Unmarshal(text, &ps)
}

// StdPermission only checks if the Ids are fully matching.
type StdPermission struct {
	IdStr string
}

func NewStdPermission(id string) *StdPermission {
	return &StdPermission{id}
}

func (p *StdPermission) Id() string {
	return p.IdStr
}

func (p *StdPermission) Match(a Permission) bool {
	return p.IdStr == a.Id()
}

func (p *StdPermission) MarshalText() (text []byte, err error) {
	var buf bytes.Buffer
	if _, err := buf.WriteRune('"'); err != nil {
		return nil, err
	}
	if _, err := buf.WriteString(p.IdStr); err != nil {
		return nil, err
	}
	if _, err := buf.WriteRune('"'); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *StdPermission) UnmarshalText(text []byte) error {
	p.IdStr = string(bytes.Trim(text, "\""))
	return nil
}

// LayerPermission firstly checks the Id of permission.
// If the Id is matched, it can be consIdered having the permission.
// Otherwise, it checks every layers of permission.
// A role which has an upper layer granted, will be granted sub-layers permissions.
type LayerPermission struct {
	IdStr string `json:"id"`
	Sep   string `json:"sep"`
}

func NewLayerPermission(id, sep string) *LayerPermission {
	return &LayerPermission{id, sep}
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

func (p *LayerPermission) MarshalText() (text []byte, err error) {
	return json.Marshal(p)
}

func (p *LayerPermission) UnmarshalText(text []byte) error {
	if err := json.Unmarshal(text, &p); err != nil {
		return err
	}
	return nil
}
