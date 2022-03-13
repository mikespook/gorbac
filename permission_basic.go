package gorbac

func NewBasicPermission[T comparable](id T) Permission[T] {
	return BasicPermission[T]{id}
}

// BasicPermission has `SID` is shorten for `serialisable ID`
type BasicPermission[T comparable] struct {
	SID T `json:"id"`
}

// ID returns id
func (p BasicPermission[T]) ID() T {
	return p.SID
}

// Match another permission
func (p BasicPermission[T]) Match(a Permission[T]) bool {
	return p.SID == a.ID()
}
