package core

import (
	"encoding/json"
	"sync"
	"time"
)

// AuditableContext holds all information captured during a request
// for generating structured audit events. It is safe for concurrent
// use via its internal mutex.
type AuditableContext struct {
	TraceID      any                    // Unique identifier for the audit trace
	UserID       any                    // Identifier for the acting user
	Action       *ActionType            // Type of action performed (e.g., Create, Update)
	Resource     *EntityType            // Resource/entity being acted upon
	ResourceID   any                    // Identifier of the resource instance
	OldData      interface{}            // Snapshot of data before the action
	NewData      interface{}            // Snapshot of data after the action
	Metadata     map[string]interface{} // Arbitrary key-value metadata
	IPAddress    string                 // Client IP address
	UserAgent    string                 // Client user-agent string
	RequestURI   string                 // HTTP request URI
	Method       string                 // HTTP method (GET, POST, etc.)
	ResponseCode int                    // HTTP response status code
	StartTime    time.Time              // Timestamp when request processing began
	EndTime      time.Time              // Timestamp when request processing completed
	mu           sync.RWMutex           // Protects mutable fields
}

// SetUserID safely sets the UserID on the context.
func (ac *AuditableContext) SetUserID(userID any) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.UserID = userID
}

// SetResourceID safely sets the ResourceID on the context.
func (ac *AuditableContext) SetResourceID(resourceID any) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.ResourceID = resourceID
}

// SetOldData safely records the pre-action data snapshot.
func (ac *AuditableContext) SetOldData(data interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.OldData = data
}

// SetNewData safely records the post-action data snapshot.
func (ac *AuditableContext) SetNewData(data interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.NewData = data
}

// SetMetadata safely sets a single metadata key-value pair.
func (ac *AuditableContext) SetMetadata(key string, value interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if ac.Metadata == nil {
		ac.Metadata = make(map[string]interface{})
	}
	ac.Metadata[key] = value
}

// SetBulkMetadata safely merges multiple metadata entries.
func (ac *AuditableContext) SetBulkMetadata(metadata map[string]interface{}) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	if ac.Metadata == nil {
		ac.Metadata = make(map[string]interface{})
	}
	for k, v := range metadata {
		ac.Metadata[k] = v
	}
}

// ToEvent converts the AuditableContext into an AuditEvent struct,
// serializing OldData and NewData to JSON if present. The returned
// AuditEvent is a snapshot and safe for publishing.
func (ac *AuditableContext) ToEvent() AuditEvent {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	event := AuditEvent{
		TraceID:      ac.TraceID,
		UserID:       ac.UserID,
		Action:       ac.Action,
		Resource:     ac.Resource,
		ResourceID:   ac.ResourceID,
		Metadata:     ac.Metadata,
		IPAddress:    ac.IPAddress,
		UserAgent:    ac.UserAgent,
		RequestURI:   ac.RequestURI,
		Method:       ac.Method,
		ResponseCode: ac.ResponseCode,
		StartTime:    ac.StartTime,
		EndTime:      ac.EndTime,
		Success:      ac.ResponseCode >= 200 && ac.ResponseCode < 400,
	}

	if ac.OldData != nil {
		if data, err := json.Marshal(ac.OldData); err == nil {
			event.OldData = data
		}
	}

	if ac.NewData != nil {
		if data, err := json.Marshal(ac.NewData); err == nil {
			event.NewData = data
		}
	}

	return event
}
