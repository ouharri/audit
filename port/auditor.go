package port

import (
	"context"
	"github.com/ouharri/audit/core"
)

// Auditor represents the component responsible for processing or transporting
// completed audit events (core.AuditEvent) to an external system.
// Implementations may publish to message queues, databases, or logging services.
type Auditor interface {
	// Audit takes a finalized AuditEvent and delivers it to the configured sink.
	// Ctx carries request-scoped values such as trace identifiers.
	Audit(ctx context.Context, event core.AuditEvent)
}
