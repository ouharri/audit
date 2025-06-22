package core

// ContextOption defines a function that applies a mutation to an AuditableContext.
// Options are typically composed and executed via SetContext.
type ContextOption func(*AuditableContext)

// WithUserID returns an option to set the user identifier on the context.
// userID may be any type representing the authenticated actor (e.g., UUID, int, string).
func WithUserID(userID any) ContextOption {
	return func(ac *AuditableContext) {
		ac.SetUserID(userID)
	}
}

// WithResourceID returns an option to set the target resource identifier.
// resourceID identifies the specific entity instance being acted upon.
func WithResourceID(resourceID any) ContextOption {
	return func(ac *AuditableContext) {
		ac.SetResourceID(resourceID)
	}
}

// WithOldData returns an option to record the pre-action data snapshot.
// oldData can be any serializable type; it will be marshaled when ToEvent is called.
func WithOldData(data interface{}) ContextOption {
	return func(ac *AuditableContext) {
		ac.SetOldData(data)
	}
}

// WithNewData returns an option to record the post-action data snapshot.
// newData can be any serializable type; it will be marshaled when ToEvent is called.
func WithNewData(data interface{}) ContextOption {
	return func(ac *AuditableContext) {
		ac.SetNewData(data)
	}
}

// WithMetadata returns an option to add a single key-value pair to the context metadata.
// Use for ad-hoc fields that supplement the event beyond CRUD data.
func WithMetadata(key string, value interface{}) ContextOption {
	return func(ac *AuditableContext) {
		ac.SetMetadata(key, value)
	}
}

// WithBulkMetadata returns an option to merge multiple key-value pairs into the context metadata.
// Useful for injecting pagination info, request parameters, or other structured data.
func WithBulkMetadata(metadata map[string]interface{}) ContextOption {
	return func(ac *AuditableContext) {
		ac.SetBulkMetadata(metadata)
	}
}
