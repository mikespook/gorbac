package gorbac

// Permission interface
// T is the type of permission ID
type Permission[T comparable] interface {
	ID() T
	Match(Permission[T]) bool
}

// Permissions list
type Permissions[T comparable] map[T]Permission[T]

func NewPermission[T comparable](id T) Permission[T] {
	return StdPermission[T]{id}
}

// BasicPermission has `SID` is shorten for `serialisable ID`
type StdPermission[T comparable] struct {
	SID T `json:"id"`
}

// ID returns id
func (p StdPermission[T]) ID() T {
	return p.SID
}

// Match another permission
func (p StdPermission[T]) Match(a Permission[T]) bool {
	return p.SID == a.ID()
}
