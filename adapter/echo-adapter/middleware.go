package echoadapter

import (
	"github.com/ouharri/audit/core"
	"github.com/ouharri/audit/transport"
	"time"

	"github.com/labstack/echo/v4"
)

// EchoMw is the Echo-specific middleware that manages audit context lifecycle
// and delegate event publication to the configured Publisher.
type EchoMw struct {
	cfg transport.Config
}

func NewEchoMiddleware(cfg transport.Config) AuditableEchoMiddleware {
	return &EchoMw{cfg}
}

// Root returns a global Echo middleware function that initializes auditing for each request.
// It records timing, user-agent, IP, URI, and response status, and triggers asynchronous publishing
// if both Resource and Action have been set in the context.
func (em *EchoMw) Root() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			if em.cfg.Skipper != nil && em.cfg.Skipper(ctx) {
				return next(c)
			}

			auditCtx := &core.AuditableContext{
				TraceID:    em.cfg.NewTraceID(),
				Metadata:   make(map[string]interface{}),
				IPAddress:  c.RealIP(),
				UserAgent:  c.Request().UserAgent(),
				RequestURI: c.Request().RequestURI,
				Method:     c.Request().Method,
				StartTime:  time.Now(),
			}

			core.SetEchoAuditContext(c, auditCtx)

			err := next(c)

			auditCtx.EndTime = time.Now()

			if resp := c.Response(); resp != nil {
				auditCtx.ResponseCode = resp.Status
			}

			if auditCtx.Resource != nil && auditCtx.Action != nil {
				go em.cfg.Auditor.Audit(ctx, auditCtx.ToEvent())
			} else {
				// If no resource or action is set, we don't send an audit event
				// This can happen if the middleware is used without specifying a resource/action
				// or if the request doesn't require auditing.
				//TODO: ?? log this case
			}

			return err
		}
	}
}

// For returns a per-route middleware factory that sets the EntityType and ActionType
// on the current AuditableContext before handler execution. It also populates UserID
// if a UserFromContext extractor is configured.
func (em *EchoMw) For(resource core.EntityType) AuditableEchoActionFactory {
	return func(action core.ActionType) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				if raw := c.Get(string(core.AuditableCtxKey)); raw != nil {
					if auditCtx, ok := raw.(*core.AuditableContext); ok {
						auditCtx.Resource = &resource
						auditCtx.Action = &action
						if em.cfg.UserFromContext != nil {
							if u := em.cfg.UserFromContext(c.Request().Context()); u != nil {
								auditCtx.UserID = u
							}
						}
					}
				}
				return next(c)
			}
		}
	}
}
