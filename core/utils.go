package core

import (
	"context"

	"github.com/labstack/echo/v4"
)

// SetContext applies the provided ContextOption functions to the AuditableContext
// stored in ctx. If no AuditableContext is present, this is a no-op.
//
// Example:
//
//	core.SetContext(ctx, core.WithResourceID(id), core.WithNewData(obj))
func SetContext(ctx context.Context, opts ...ContextOption) {
	if ac := GetAuditContext(ctx); ac != nil {
		for _, opt := range opts {
			opt(ac)
		}
	}
}

// SetEchoAuditContext stores the given AuditableContext in both Echo's
// internal context and the underlying http.Request.Context, ensuring
// later handlers and middleware can retrieve it via GetAuditContext.
//
// Example (used by Echo adapter):
//
//	core.SetEchoAuditContext(c, auditCtx)
func SetEchoAuditContext(c echo.Context, auditCtx *AuditableContext) {
	// Store in Echo's context map
	c.Set(string(AuditableCtxKey), auditCtx)

	// Also inject into the request's Context
	req := c.Request()
	ctxWithAudit := context.WithValue(req.Context(), AuditableCtxKey, auditCtx)
	c.SetRequest(req.WithContext(ctxWithAudit))
}

// GetAuditContext retrieves the AuditableContext from the provided context.Context.
// Returns nil if no AuditableContext is stored.
//
// Use this within business logic or decorators to mutate audit data:
//
//	ac := core.GetAuditContext(ctx)
func GetAuditContext(ctx context.Context) *AuditableContext {
	if raw := ctx.Value(AuditableCtxKey); raw != nil {
		if ac, ok := raw.(*AuditableContext); ok {
			return ac
		}
	}
	return nil
}
