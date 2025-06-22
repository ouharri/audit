package port

import "github.com/ouharri/audit/core"

// AuditableActionFactory produces a middleware handler of type T
// that injects a specific ActionType into the existing AuditableContext for a given route.
// T represents the middleware function type for the target framework.
type AuditableActionFactory[T any] func(action core.ActionType) T

// Middleware abstracts a transport middleware entry point (Root)
// and a per-resource factory (For). T represents the middleware function type
// for the target framework (e.g., echo.MiddlewareFunc, gin.HandlerFunc).
type Middleware[T any] interface {
	// Root returns the global middleware function that initializes audit context.
	Root() T

	// For returns a factory that attaches the given EntityType and ActionType
	// to the audit context for a route or handler.
	For(resource core.EntityType) AuditableActionFactory[T]
}
