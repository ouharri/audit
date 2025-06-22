package core

import (
	"encoding/json"
	"time"
)

// auditCtxKey is the private type used as the context key for storing
// the AuditableContext in a request's Context.
type auditCtxKey string

// AuditableCtxKey is the key under which the AuditableContext is stored
// in context.Context and Echo Context values.
var AuditableCtxKey = auditCtxKey("_audit_ctx")

// AuditEvent represents a fully materialized audit record ready for serialization
// and dispatch to an audit event Publisher. All timestamps, identifiers, and
// serialized data snapshots are captured within this struct.
type AuditEvent struct {
	TraceID      any                    `json:"traceId"`
	UserID       any                    `json:"userId,omitempty"`
	Action       *ActionType            `json:"action,omitempty"`
	Resource     *EntityType            `json:"resource,omitempty"`
	ResourceID   any                    `json:"resourceId,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	IPAddress    string                 `json:"ipAddress,omitempty"`
	UserAgent    string                 `json:"userAgent,omitempty"`
	RequestURI   string                 `json:"requestUri,omitempty"`
	Method       string                 `json:"method,omitempty"`
	ResponseCode int                    `json:"responseCode,omitempty"`
	Success      bool                   `json:"success,omitempty"`
	StartTime    time.Time              `json:"startTime"`
	EndTime      time.Time              `json:"endTime"`
	OldData      json.RawMessage        `json:"oldData,omitempty"`
	NewData      json.RawMessage        `json:"newData,omitempty"`
}

// ActionType defines the type of operation being audited.
// Typical values include CRUD actions (e.g., "Create", "Update").
type ActionType string

// EntityType defines the resource or domain entity being audited.
// Typical values correspond to application domain types (e.g., "User", "Order").
type EntityType string
