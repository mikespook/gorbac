package gorbac

// Permission interface
// T is the type of permission ID
type Permission[T comparable] interface {
	ID() T
	Match(Permission[T]) bool
}

// Permissions list
type Permissions[T comparable] map[T]Permission[T]
