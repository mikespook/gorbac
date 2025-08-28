# goRBAC - LLM Coder Guide

This document provides a comprehensive guide to the goRBAC (Role-Based Access Control) package for LLM coders who need to understand and work with this Go library without going through the entire codebase.

## Overview

goRBAC is a lightweight role-based access control implementation in Golang. It provides a simple and efficient way to manage roles, permissions, and their relationships in applications that require access control.

### Core Concepts

1. **Identity**: An entity that has one or more roles.
2. **Role**: A named entity that can be assigned permissions and can inherit from parent roles.
3. **Permission**: An entity that represents an action or resource access right.
4. **Inheritance**: Roles can inherit permissions from parent roles, forming a hierarchical structure.

### Key Features

- Generic support for different ID types (string, int, etc.)
- Role inheritance with circular dependency detection
- Thread-safe operations
- JSON serialization support
- Extensible interfaces for custom implementations
- Built-in utility functions for common operations

## Package Structure

```
gorbac/
├── rbac.go              # Main RBAC implementation
├── role.go              # Role interface and standard implementation
├── permission.go        # Permission interface and standard implementation
├── helper.go            # Utility functions
├── helper_test.go       # Tests for helper functions
├── rbac_test.go         # Tests for RBAC implementation
├── role_test.go         # Tests for role implementation
├── permission_test.go   # Tests for permission implementation
├── example_test.go      # Usage examples
├── examples/            # Complete example applications
│   ├── persistence/     # Example showing data persistence
│   └── user-defined/    # Example with custom role implementation
├── README.md            # Project documentation
└── go.mod               # Go module definition
```

## Core Components

### 1. RBAC Structure (`rbac.go`)

The main RBAC structure manages roles and their inheritance relationships.

#### Key Methods

- `New[T comparable]() *RBAC[T]` - Creates a new RBAC instance
- `Add(r Role[T]) error` - Adds a role to the RBAC instance
- `Remove(id T) error` - Removes a role by ID
- `Get(id T) (Role[T], []T, error)` - Gets a role and its parents
- `SetParent(id T, parent T) error` - Sets a parent for a role
- `SetParents(id T, parents []T) error` - Sets multiple parents for a role
- `GetParents(id T) ([]T, error)` - Gets all parents of a role
- `RemoveParent(id T, parent T) error` - Removes a parent from a role
- `IsGranted(id T, p Permission[T], assert AssertionFunc[T]) bool` - Checks if a role has a permission

#### Thread Safety

All operations on the RBAC structure are thread-safe using read-write mutexes.

### 2. Role Implementation (`role.go`)

The `Role[T]` struct is the default implementation:

```go
type Role[T comparable] struct {
    sync.RWMutex
    ID          T `json:"id"`
    permissions Permissions[T]
}
```

#### Key Methods

- `NewRole[T comparable](id T) Role[T]` - Creates a new role
- `Assign(p Permission[T]) error` - Assigns a permission to the role
- `Permit(p Permission[T]) bool` - Checks if the role has a specific permission
- `Revoke(p Permission[T]) error` - Revokes a permission from the role
- `Permissions() []Permission[T]` - Returns all permissions assigned to the role

### 3. Permission Interface and Implementation (`permission.go`)

The `Permission[T]` interface defines the contract for permissions:

```go
type Permission[T comparable] interface {
    ID() T
    Match(Permission[T]) bool
}
```

#### Standard Permission Implementation

The package provides `StdPermission[T]` as the default implementation:

- `SID` - Serializable ID of the permission

#### Key Methods

- `NewPermission[T comparable](id T) Permission[T]` - Creates a new permission
- `ID() T` - Returns the permission ID
- `Match(Permission[T]) bool` - Checks if this permission matches another

### 4. Helper Functions (`helper.go`)

Utility functions for common operations:

#### Walk Function

- `Walk[T comparable](rbac *RBAC[T], h WalkHandler[T]) error` - Iterates through all roles

#### Inheritance Validation

- `InherCircle[T comparable](rbac *RBAC[T]) error` - Detects circular inheritance

#### Permission Checking

- `AnyGranted[T comparable](rbac *RBAC[T], roles []T, permission Permission[T], assert AssertionFunc[T]) bool` - Checks if any role has a permission
- `AllGranted[T comparable](rbac *RBAC[T], roles []T, permission Permission[T], assert AssertionFunc[T]) bool` - Checks if all roles have a permission

## Usage Examples

### Basic Usage

```go
// Create a new RBAC instance
rbac := gorbac.New[string]()

// Create roles
rA := gorbac.NewRole("role-a")
rB := gorbac.NewRole("role-b")

// Create permissions
pA := gorbac.NewPermission("permission-a")
pB := gorbac.NewPermission("permission-b")

// Assign permissions to roles
rA.Assign(pA)
rB.Assign(pB)

// Add roles to RBAC
rbac.Add(rA)
rbac.Add(rB)

// Set inheritance
rbac.SetParent("role-a", "role-b")

// Check permissions
if rbac.IsGranted("role-a", pA, nil) {
    // role-a has permission-a
}
```

### Working with Different ID Types

The package supports generic ID types:

```go
// String IDs
rbacStr := gorbac.New[string]()

// Integer IDs
rbacInt := gorbac.New[int]()

// Custom struct IDs
type RoleID struct {
    Name string
    Type string
}
rbacStruct := gorbac.New[RoleID]()
```

### Custom Assertion Functions

You can provide custom assertion functions for fine-grained control:

```go
assertFunc := func(r *gorbac.RBAC[string], id string, p gorbac.Permission[string]) bool {
    // Custom logic to determine if permission should be granted
    return true // or false
}

if rbac.IsGranted("role-a", pA, assertFunc) {
    // Permission granted based on custom logic
}
```

## Persistence

The package doesn't include built-in persistence but provides mechanisms for implementing it:

### Example Persistence Approach

See `examples/persistence/persistence.go` for a complete example of:

1. Loading roles and permissions from JSON files
2. Building the RBAC structure from persisted data
3. Saving the RBAC structure back to JSON files

### Key Concepts for Persistence

1. Serialize roles and their permissions
2. Serialize inheritance relationships
3. Reconstruct the RBAC instance from persisted data

## Custom Implementations

### Custom Role Implementation

You can create custom roles by embedding the standard role:

```go
type myRole struct {
    gorbac.Role[string]  // Embed the standard role
    Label       string
    Description string
}
```

### Custom Permission Implementation

You can implement the `Permission[T]` interface to create custom permissions with additional logic in the `Match` method.

## Error Handling

The package defines standard errors:

- `ErrRoleNotExist` - When a role doesn't exist
- `ErrRoleExist` - When trying to add a role that already exists
- `ErrFoundCircle` - When circular inheritance is detected

Always check and handle these errors appropriately in your applications.

## Performance Considerations

- RBAC operations use read-write mutexes for thread safety
- Permission checking with inheritance uses recursive traversal
- Circular inheritance detection uses depth-first search
- Consider caching results for frequently checked permissions in performance-critical applications

## Testing

The package includes comprehensive tests covering:

- Basic RBAC operations
- Role and permission management
- Inheritance relationships
- Circular dependency detection
- Helper functions
- Various ID types

See the `*_test.go` files for detailed usage examples.

## Quick Reference

| Component | File | Key Functions |
|-----------|------|---------------|
| RBAC Core | `rbac.go` | `New`, `Add`, `Remove`, `IsGranted`, `SetParent` |
| Roles | `role.go` | `NewRole`, `Assign`, `Permit`, `Revoke` |
| Permissions | `permission.go` | `NewPermission`, `Match` |
| Utilities | `helper.go` | `Walk`, `InherCircle`, `AnyGranted`, `AllGranted` |
| Examples | `example_test.go` | Complete usage examples |

## Common Patterns

### 1. Initialization Pattern

```go
rbac := gorbac.New[string]()
// Create roles and permissions
// Assign permissions to roles
// Add roles to RBAC
// Set up inheritance
```

### 2. Permission Checking Pattern

```go
if rbac.IsGranted("user-role", requiredPermission, nil) {
    // Allow access
} else {
    // Deny access
}
```

### 3. Batch Permission Checking

```go
roles := []string{"role1", "role2", "role3"}
if gorbac.AnyGranted(rbac, roles, permission, nil) {
    // At least one role has the permission
}

if gorbac.AllGranted(rbac, roles, permission, nil) {
    // All roles have the permission
}
```

## Extending the Package

1. Embed standard `Role[T]` struct for domain-specific role behavior
2. Implement custom `Permission[T]` interfaces for complex permission matching logic
3. Use the `Walk` function to export RBAC state for persistence
4. Add middleware functions for logging or metrics around RBAC operations

This guide provides a comprehensive overview of the goRBAC package. For implementation details, refer to the source files in the package structure.