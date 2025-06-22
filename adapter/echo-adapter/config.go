package echoadapter

import (
	"github.com/ouharri/audit/core"
	"github.com/ouharri/audit/transport"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	once sync.Once

	// mw holds the singleton Echo middleware instance.
	mw AuditableEchoMiddleware
)

// Configure wires your core audit.Auditor into the Echo middleware singleton.
// Call this once at startup.
func Configure(cfg transport.Config) {
	once.Do(func() {
		mw = NewEchoMiddleware(cfg)
	})
}

// Root returns the global Echo middleware. Panics if Configure wasn’t called.
func Root() echo.MiddlewareFunc {
	if mw == nil {
		panic("echo audit: call Configure(auditor) before Root()")
	}
	return mw.Root()
}

// For returns the per-route factory. Panics if Configure wasn’t called.
func For(resource core.EntityType) AuditableEchoActionFactory {
	if mw == nil {
		panic("echo audit: call Configure(auditor) before For()")
	}
	return mw.For(resource)
}
