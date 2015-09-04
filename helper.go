package gorbac

import "fmt"

// Return an error when detects any circle inheritance.
func InherCircle(rbac *Rbac) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()

	skipped := make(map[string]struct{})
	stack := make([]string, 0)

	for _, role := range rbac.roles {
		if err := dfs(rbac, role, skipped, stack); err != nil {
			return err
		}
	}
	return nil
}

func dfs(rbac *Rbac, role Role, skipped map[string]struct{}, stack []string) error {
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
