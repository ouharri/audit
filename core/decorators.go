package core

import (
	"context"
	"encoding/json"
)

// AuditableCreate records a "create" operation by setting the NewData field
// on the AuditableContext associated with ctx.
func AuditableCreate(ctx context.Context, resourceID any, newData interface{}) {
	SetContext(ctx,
		WithResourceID(resourceID),
		WithNewData(newData),
	)
}

// AuditableUpdate records an "update" operation by setting ResourceID,
// OldData, and NewData on the AuditableContext associated with ctx.
func AuditableUpdate(ctx context.Context, resourceID any, oldData, newData interface{}) {
	SetContext(ctx,
		WithResourceID(resourceID),
		WithOldData(oldData),
		WithNewData(newData),
	)
}

// AuditableDelete records a "delete" operation by setting ResourceID and
// OldData on the AuditableContext associated with ctx.
func AuditableDelete(ctx context.Context, resourceID any, oldData interface{}) {
	SetContext(ctx,
		WithResourceID(resourceID),
		WithOldData(oldData),
	)
}

// AuditableGet records a "get" (read) operation by setting ResourceID on
// the AuditableContext associated with ctx.
func AuditableGet(ctx context.Context, resourceID any) {
	SetContext(ctx, WithResourceID(resourceID))
}

// AuditableList records a listing operation by merging the provided metadata
// map into the AuditableContext associated with ctx.
func AuditableList(ctx context.Context, metadata map[string]interface{}) {
	if metadata != nil {
		SetContext(ctx, WithBulkMetadata(metadata))
	}
}

// AuditablePage records a paginated listing operation by serializing the
// page data into a metadata map and merging it into the context.
// If serialization fails, no metadata is set.
func AuditablePage(ctx context.Context, pageData interface{}) {
	bytes, err := json.Marshal(pageData)
	if err != nil {
		// Serialization failed: skip setting page metadata
		return
	}
	var metadata map[string]interface{}
	if err := json.Unmarshal(bytes, &metadata); err != nil {
		// Unmarshal failed: skip setting page metadata
		return
	}
	SetContext(ctx, WithBulkMetadata(metadata))
}

// AuditableAction records a custom action by setting ResourceID and
// merging arbitrary metadata into the AuditableContext associated with ctx.
func AuditableAction(ctx context.Context, resourceID any, metadata map[string]interface{}) {
	opts := []ContextOption{WithResourceID(resourceID)}
	if metadata != nil {
		opts = append(opts, WithBulkMetadata(metadata))
	}
	SetContext(ctx, opts...)
}
