package gorbac

import "fmt"

// InherCircle returns an error when detecting any circle inheritance.
func InherCircle(rbac *RBAC) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()

	skipped := make(map[string]struct{})
	var stack []string

	for _, role := range rbac.roles {
		if err := dfs(rbac, role, skipped, stack); err != nil {
			return err
		}
	}
	return nil
}

func dfs(rbac *RBAC, role Role, skipped map[string]struct{}, stack []string) error {
	if _, ok := skipped[role.Name()]; ok {
		return nil
	}
	for _, item := range stack {
		if item == role.Name() {
			return fmt.Errorf("Found circle: %s", stack)
		}
	}
	if len(role.Parents()) == 0 {
		stack = make([]string, 0)
		skipped[role.Name()] = struct{}{}
		return nil
	}
	stack = append(stack, role.Name())
	for _, pname := range role.Parents() {
		if err := dfs(rbac, rbac.roles[pname], skipped, stack); err != nil {
			return err
		}
	}
	return nil
}

// AnyGranted checks if any role has the permission.
func AnyGranted(rbac *RBAC, roles []string, permission Permission,
	assert AssertionFunc) bool {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	for _, role := range roles {
		if rbac.isGranted(role, permission, assert) {
			return true
		}
	}
	return false
}

// AllGranted checks if all roles have the permission.
func AllGranted(rbac *RBAC, roles []string, permission Permission,
	assert AssertionFunc) bool {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	for _, role := range roles {
		if !rbac.isGranted(role, permission, assert) {
			return false
		}
	}
	return true
}
