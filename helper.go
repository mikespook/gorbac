package gorbac

import "fmt"

// WalkHandler is a function defined by user to handle role
type WalkHandler[T comparable] func(Role[T], []T) error

// Walk passes each Role to WalkHandler
func Walk[T comparable](rbac *RBAC[T], h WalkHandler[T]) (err error) {
	if h == nil {
		return
	}
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	for id := range rbac.roles {
		var parents []T
		r := rbac.roles[id]
		for parent := range rbac.parents[id] {
			parents = append(parents, parent)
		}
		if err := h(r, parents); err != nil {
			return err
		}
	}
	return
}

// InherCircle returns an error when detecting any circle inheritance.
func InherCircle[T comparable](rbac *RBAC[T]) (err error) {
	rbac.mutex.Lock()

	skipped := make(map[T]struct{}, len(rbac.roles))
	var stack []T

	for id := range rbac.roles {
		if err = dfs(rbac, id, skipped, stack); err != nil {
			break
		}
	}
	rbac.mutex.Unlock()
	return err
}

var (
	ErrFoundCircle = fmt.Errorf("Found circle")
)

// https://en.wikipedia.org/wiki/Depth-first_search
func dfs[T comparable](rbac *RBAC[T], id T, skipped map[T]struct{},
	stack []T) error {
	if _, ok := skipped[id]; ok {
		return nil
	}
	for _, item := range stack {
		if item == id {
			return ErrFoundCircle
		}
	}
	parents := rbac.parents[id]
	if len(parents) == 0 {
		stack = nil
		skipped[id] = empty
		return nil
	}
	stack = append(stack, id)
	for pid := range parents {
		if err := dfs(rbac, pid, skipped, stack); err != nil {
			return err
		}
	}
	return nil
}

// AnyGranted checks if any role has the permission.
func AnyGranted[T comparable](rbac *RBAC[T], roles []T,
	permission Permission[T], assert AssertionFunc[T]) (ok bool) {
	rbac.mutex.Lock()
	for _, role := range roles {
		if rbac.isGranted(role, permission, assert) {
			ok = true
			break
		}
	}
	rbac.mutex.Unlock()
	return
}

// AllGranted checks if all roles have the permission.
func AllGranted[T comparable](rbac *RBAC[T], roles []T,
	permission Permission[T], assert AssertionFunc[T]) (ok bool) {
	ok = true
	rbac.mutex.Lock()
	for _, role := range roles {
		if !rbac.isGranted(role, permission, assert) {
			ok = false
			break
		}
	}
	rbac.mutex.Unlock()
	return
}
