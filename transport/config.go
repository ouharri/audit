package transport

import (
	"context"
	"github.com/ouharri/audit/port"
)

// Config encapsulates the components required by an HTTP transport adapter
// for audit logging. Populate this struct with your Publisher (Auditor),
// trace ID generator, user extractor, and optional request-skipping logic.
//
// Example:
//
//	cfg := common.Config{
//	    // Required: your audit event publisher
//	    Auditor: myPublisher,
//
//	    // Required: generate a new trace ID per request
//	    NewTraceID: func() any { return uuid.New() },
//
//	    // Optional: extract the current user ID from context
//	    UserFromContext: extractUserID,
//
//	    // Optional: skip auditing for health checks or other routes
//	    Skipper: func(ctx context.Context) bool {
//	        return skipHealthCheck(ctx)
//	    },
//	}
type Config struct {
	// Auditor processes and delivers audit events. Required and must implement port.Auditor.
	Auditor port.Auditor

	// NewTraceID generates a new trace identifier for each request.
	// Typically implemented using a UUID generator or a distributed
	// tracing library. Its return value is stored as TraceID in events.
	NewTraceID func() any

	// UserFromContext extracts the current user (string, int, uuid, etc.)
	UserFromContext func(ctx context.Context) any

	// Skipper determines whether to skip auditing for a given request.
	// If nil, auditing is always enabled.
	Skipper func(ctx context.Context) bool
}
