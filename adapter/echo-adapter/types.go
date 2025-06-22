package echoadapter

import (
	"github.com/labstack/echo/v4"
	"github.com/ouharri/audit/port"
)

// AuditableEchoActionFactory is the Echo-specific middleware factory type.
// It injects an ActionType into the audit context for a route.
// Equivalent to port.AuditableActionFactory[echo.MiddlewareFunc].
type AuditableEchoActionFactory = port.AuditableActionFactory[echo.MiddlewareFunc]

// AuditableEchoMiddleware is the Echo-specific middleware interface.
// It specializes the generic port.Middleware[T] for Echo.
type AuditableEchoMiddleware = port.Middleware[echo.MiddlewareFunc]
