package gorbac

import "fmt"

// WalkHandler is a function defined by user to handle role
type WalkHandler[K comparable] func(Role[K], []K) error

// Walk passes each Role to WalkHandler
func Walk[K comparable](rbac *RBAC[K], h WalkHandler[K]) (err error) {
	if h == nil {
		return
	}
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	for id := range rbac.roles {
		var parents []K
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
func InherCircle[K comparable](rbac *RBAC[K]) (err error) {
	rbac.mutex.Lock()

	skipped := make(map[K]struct{}, len(rbac.roles))
	var stack []K

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
func dfs[K comparable](rbac *RBAC[K], id K, skipped map[K]struct{}, stack []K) error {
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
func AnyGranted[K comparable](rbac *RBAC[K], roles []K, permission Permission[K],
	assert AssertionFunc[K]) (ok bool) {
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
func AllGranted[K comparable](rbac *RBAC[K], roles []K, permission Permission[K],
	assert AssertionFunc[K]) (ok bool) {
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
