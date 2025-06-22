# Audit Suite

_Empowering Trust Through Seamless, Secure Auditing for Go Applications_

[![last-commit](https://img.shields.io/github/last-commit/ouharri/audit?style=flat&logo=git&logoColor=white&color=0080ff)](https://github.com/ouharri/audit/commits/main)
[![repo-top-language](https://img.shields.io/github/languages/top/ouharri/audit?style=flat&color=0080ff)](https://github.com/ouharri/audit)
[![repo-language-count](https://img.shields.io/github/languages/count/ouharri/audit?style=flat&color=0080ff)](https://github.com/ouharri/audit)
[![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=flat&logo=Go&logoColor=white)](https://golang.org/)
[![License](https://img.shields.io/github/license/ouharri/audit?color=0080ff)](LICENSE)

---

## Table of Contents

- [Overview](#overview)
- [github.com/ouharri/audit (Core Library)](#githubcomouharraudit-core-library)
    - [Overview](#core-overview)
    - [Project Layout](#core-project-layout)
    - [Installation](#core-installation)
    - [Configuration](#core-configuration)
    - [Core API Reference](#core-api-reference)
        - [Transport Configuration](#transport-configuration)
        - [Port Interfaces](#port-interfaces)
        - [Core Domain Types](#core-domain-types)
        - [Context Management](#context-management)
        - [Audit Decorators](#audit-decorators)
        - [Context Options](#context-options)
    - [Usage Examples](#core-usage-examples)
    - [Echo Adapter](#echo-adapter)
        - [Echo Configuration](#echo-configuration)
        - [Echo API Reference](#echo-api-reference)
        - [Echo Usage Examples](#echo-usage-examples)
- [Complete Integration Example](#complete-integration-example)
- [Best Practices](#best-practices)
- [License](#license)

---

## Overview

**Audit Suite** is a comprehensive, modular Go library for implementing robust audit logging in your applications. It provides a clean, extensible architecture for capturing, enriching, and delivering audit events with support for popular web frameworks through dedicated adapters.

**Key Features:**
- 🏗️ **Clean Architecture**: Separation of concerns with ports and adapters pattern
- 🔧 **Framework Agnostic**: Core library works with any Go application
- 📊 **Rich Event Model**: Comprehensive audit events with metadata, timing, and data snapshots
- 🔒 **Secure Context**: Thread-safe context management with mutex protection
- 🚀 **Pluggable Backends**: Interface-driven design for any audit destination
- 🌐 **Echo Integration**: Ready-to-use middleware for Echo web framework
- 📈 **Performance Focused**: Asynchronous event publishing with minimal overhead

---

## github.com/ouharri/audit (Core Library)

### Core Overview

The `github.com/ouharri/audit` core library provides the foundational components for audit logging in Go applications. It implements a clean architecture with clear separation between domain logic, ports (interfaces), and adapters (implementations).

**Module Path:** `github.com/ouharri/audit`

#### Use Cases

- **Compliance Auditing**: Track all CRUD operations for regulatory compliance
- **Security Monitoring**: Monitor user actions and detect suspicious behavior
- **Business Intelligence**: Capture business events for analytics and reporting
- **Debugging & Tracing**: Correlate application events across microservices
- **Change Tracking**: Maintain audit trails for data modifications

---

### Core Project Layout

```
github.com/ouharri/audit/
├── core/
│   ├── context.go         # AuditableContext and context management
│   ├── decorators.go      # High-level audit decorators (Create, Update, etc.)
│   ├── domain.go          # Core domain types and AuditEvent
│   ├── options.go         # Functional options for context modification
│   └── utils.go           # Context utilities and helpers
├── port/
│   ├── auditor.go         # Auditor interface for event publishing
│   └── middleware.go      # Generic middleware interfaces
├── transport/
│   └── config.go          # Transport configuration structure
├── echoadapter/
│   ├── middleware.go      # Echo-specific middleware implementation
│   ├── singleton.go       # Singleton pattern for global configuration
│   └── types.go           # Echo-specific type aliases
├── go.mod
├── go.sum
└── LICENSE
```

---

### Core Installation

```sh
go get github.com/ouharri/audit
```

---

### Core Configuration

Configure the audit system with your specific requirements:

```go
import (
    "context"
    "github.com/ouharri/audit/transport"
    "github.com/ouharri/audit/port"
)

// Configure your audit system
config := transport.Config{
    // Required: Your audit event publisher
    Auditor: myAuditor, // implements port.Auditor

    // Required: Generate unique trace IDs
    NewTraceID: func() any { 
        return uuid.New().String() 
    },

    // Optional: Extract user from context
    UserFromContext: func(ctx context.Context) any {
        if userID := ctx.Value("user_id"); userID != nil {
            return userID
        }
        return nil
    },

    // Optional: Skip auditing for certain requests
    Skipper: func(ctx context.Context) bool {
        // Skip health checks, metrics endpoints, etc.
        return isHealthCheck(ctx)
    },
}
```

---

### Core API Reference

#### Transport Configuration

```go
// Config encapsulates all components required by transport adapters
type Config struct {
    // Required: Processes and delivers audit events
    Auditor port.Auditor

    // Required: Generates unique trace identifiers
    NewTraceID func() any

    // Optional: Extracts current user from context
    UserFromContext func(ctx context.Context) any

    // Optional: Determines whether to skip auditing
    Skipper func(ctx context.Context) bool
}
```

#### Port Interfaces

```go
// Auditor processes completed audit events
type Auditor interface {
    // Audit delivers an event to the configured destination
    Audit(ctx context.Context, event core.AuditEvent) error
}

// Middleware abstracts transport-specific middleware
type Middleware[T any] interface {
    // Root returns the global middleware function
    Root() T

    // For returns a factory for resource-specific middleware
    For(resource core.EntityType) AuditableActionFactory[T]
}

// AuditableActionFactory creates action-specific middleware
type AuditableActionFactory[T any] func(action core.ActionType) T
```

#### Core Domain Types

```go
// AuditEvent represents a complete audit record
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

// ActionType defines the operation being audited
type ActionType string

// EntityType defines the resource being audited
type EntityType string

// Common action types (must define these in your application)
const (
    ActionCreate ActionType = "CREATE"
    ActionRead   ActionType = "READ"
    ActionUpdate ActionType = "UPDATE"
    ActionDelete ActionType = "DELETE"
    ActionList   ActionType = "LIST"
)

// AuditableContext holds request-scoped audit information
type AuditableContext struct {
    TraceID      any                    // Unique trace identifier
    UserID       any                    // Acting user identifier
    Action       *ActionType            // Operation type
    Resource     *EntityType            // Target resource type
    ResourceID   any                    // Specific resource instance
    OldData      interface{}            // Pre-action data snapshot
    NewData      interface{}            // Post-action data snapshot
    Metadata     map[string]interface{} // Additional context
    IPAddress    string                 // Client IP address
    UserAgent    string                 // Client user agent
    RequestURI   string                 // Request URI
    Method       string                 // HTTP method
    ResponseCode int                    // Response status code
    StartTime    time.Time              // Request start time
    EndTime      time.Time              // Request end time
    // mu           sync.RWMutex        // Thread safety (unexported)
}
```

#### Context Management

```go
// SetContext applies options to the audit context in ctx
func SetContext(ctx context.Context, opts ...ContextOption)

// SetEchoAuditContext stores context in Echo framework
func SetEchoAuditContext(c echo.Context, auditCtx *AuditableContext)

// GetAuditContext retrieves audit context from ctx
func GetAuditContext(ctx context.Context) *AuditableContext

// Thread-safe context methods
func (ac *AuditableContext) SetUserID(userID any)
func (ac *AuditableContext) SetResourceID(resourceID any)
func (ac *AuditableContext) SetOldData(data interface{})
func (ac *AuditableContext) SetNewData(data interface{})
func (ac *AuditableContext) SetMetadata(key string, value interface{})
func (ac *AuditableContext) SetBulkMetadata(metadata map[string]interface{})

// ToEvent converts context to publishable event
func (ac *AuditableContext) ToEvent() *AuditEvent
```

#### Audit Decorators

High-level functions for common audit operations:

```go
// AuditableCreate records a create operation
func AuditableCreate(ctx context.Context, newData interface{})

// AuditableUpdate records an update operation
func AuditableUpdate(ctx context.Context, resourceID any, oldData, newData interface{})

// AuditableDelete records a delete operation
func AuditableDelete(ctx context.Context, resourceID any, oldData interface{})

// AuditableGet records a read operation
func AuditableGet(ctx context.Context, resourceID any)

// AuditableList records a list operation
func AuditableList(ctx context.Context, metadata map[string]interface{})

// AuditablePage records a paginated list operation
func AuditablePage(ctx context.Context, pageData interface{})

// AuditableAction records a custom action
func AuditableAction(ctx context.Context, resourceID any, metadata map[string]interface{})
```

#### Context Options

Functional options for modifying audit context:

```go
// WithUserID sets the acting user
func WithUserID(userID any) ContextOption

// WithResourceID sets the target resource instance
func WithResourceID(resourceID any) ContextOption

// WithOldData sets the pre-action data snapshot
func WithOldData(data interface{}) ContextOption

// WithNewData sets the post-action data snapshot
func WithNewData(data interface{}) ContextOption

// WithMetadata adds a single metadata key-value pair
func WithMetadata(key string, value interface{}) ContextOption

// WithBulkMetadata merges multiple metadata entries
func WithBulkMetadata(metadata map[string]interface{}) ContextOption
```

---

### Core Usage Examples

**Basic Context Manipulation:**

```go
package main

import (
    "context"
    "github.com/ouharri/audit/core"
)

func updateUser(ctx context.Context, userID string, oldUser, newUser User) error {
    // Record the update operation
    core.SetContext(ctx,
        core.WithResourceID(userID),
        core.WithOldData(oldUser),
        core.WithNewData(newUser),
    )

    // Perform your business logic
    return userService.Update(userID, newUser)
}
```

**Using Audit Decorators:**

```go
package main

import (
    "context"
    "github.com/ouharri/audit/core"
)

func createProduct(ctx context.Context, product Product) (*Product, error) {
    // Perform creation
    created, err := productService.Create(product)
    if err != nil {
        return nil, err
    }

    // Record the creation
    core.AuditableCreate(ctx, created)
    
    return created, nil
}

func deleteProduct(ctx context.Context, productID string) error {
    // Get existing data before deletion
    existing, err := productService.GetByID(productID)
    if err != nil {
        return err
    }

    // Perform deletion
    if err := productService.Delete(productID); err != nil {
        return err
    }

    // Record the deletion
    core.AuditableDelete(ctx, productID, existing)
    
    return nil
}
```

**Custom Auditor Implementation:**

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "github.com/ouharri/audit/core"
)

type LogAuditor struct {
    logger *log.Logger
}

func (la *LogAuditor) Audit(ctx context.Context, event core.AuditEvent) error {
    eventJSON, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    la.logger.Printf("AUDIT: %s", string(eventJSON))
    return nil
}

// Usage
auditor := &LogAuditor{logger: log.Default()}
```

---

### Echo Adapter

The Echo adapter provides seamless integration with the Echo web framework.

#### Echo Configuration

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/ouharri/audit/echoadapter"
    "github.com/ouharri/audit/transport"
)

func main() {
    e := echo.New()

    // Initialize Echo adapter
    echoadapter.Configure(transport.Config{
		Auditor:    myAuditor,
		NewTraceID: generateTraceID,
		UserFromContext: extractUser,
	})

    // Apply global audit middleware
    e.Use(echoadapter.Root())

    // Configure routes with specific audit settings
    setupRoutes(e)

    e.Logger.Fatal(e.Start(":8080"))
}
```

#### Echo API Reference

```go
// Configure initializes the Echo adapter with audit configuration
func Configure(cfg transport.Config)

// Root returns the global Echo middleware function
func Root() echo.MiddlewareFunc

// For returns a factory for resource-specific middleware
func For(resource core.EntityType) AuditableEchoActionFactory

// Type aliases for Echo integration
type AuditableEchoActionFactory = port.AuditableActionFactory[echo.MiddlewareFunc]
type AuditableEchoMiddleware = port.Middleware[echo.MiddlewareFunc]
```

#### Echo Usage Examples

**Basic Route Configuration:**

```go
func setupRoutes(e *echo.Echo) {
    // Define resource types
    const (
        EntityUser    core.EntityType = "USER"
        EntityProduct core.EntityType = "PRODUCT"
    )

    // Create factories for different resources
    userAudit := echoadapter.For(EntityUser)
    productAudit := echoadapter.For(EntityProduct)

    // Configure user routes
    userGroup := e.Group("/users")
    userGroup.POST("", createUserHandler, userAudit(core.ActionCreate))
    userGroup.GET("/:id", getUserHandler, userAudit(core.ActionRead))
    userGroup.PUT("/:id", updateUserHandler, userAudit(core.ActionUpdate))
    userGroup.DELETE("/:id", deleteUserHandler, userAudit(core.ActionDelete))

    // Configure product routes
    productGroup := e.Group("/products")
    productGroup.POST("", createProductHandler, productAudit(core.ActionCreate))
    productGroup.GET("", listProductsHandler, productAudit(core.ActionList))
}
```

**Handler Implementation with Audit Context:**

```go
func updateUserHandler(c echo.Context) error {
    ctx := c.Request().Context()
    userID := c.Param("id")

    var updateReq UpdateUserRequest
    if err := c.Bind(&updateReq); err != nil {
        return err
    }

    // Get existing user for audit trail
    existingUser, err := userService.GetByID(ctx, userID)
    if err != nil {
        return err
    }

    // Update user
    updatedUser, err := userService.Update(ctx, userID, updateReq)
    if err != nil {
        return err
    }

    // Record the update operation
    core.AuditableUpdate(ctx, userID, existingUser, updatedUser)

    return c.JSON(http.StatusOK, updatedUser)
}
```

---

## Complete Integration Example

Here's a complete example demonstrating the audit system in a real Echo application:

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/ouharri/audit/core"
    "github.com/ouharri/audit/echoadapter"
    "github.com/ouharri/audit/transport"
)

// Domain types
type User struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Created  time.Time `json:"created"`
    Modified time.Time `json:"modified"`
}

// Resource types
const (
    EntityUser core.EntityType = "USER"
)

// Action types
const (
    ActionCreate core.ActionType = "CREATE"
    ActionRead   core.ActionType = "READ"
    ActionUpdate core.ActionType = "UPDATE"
    ActionDelete core.ActionType = "DELETE"
    ActionList   core.ActionType = "LIST"
)

// Simple in-memory auditor
type ConsoleAuditor struct {
    logger *log.Logger
}

func (ca *ConsoleAuditor) Audit(ctx context.Context, event core.AuditEvent) error {
    eventData, _ := json.MarshalIndent(event, "", "  ")
    ca.logger.Printf("🔍 AUDIT EVENT:\n%s\n", string(eventData))
    return nil
}

// Mock user service
type UserService struct {
    users  map[int]*User
    nextID int
}

func NewUserService() *UserService {
    return &UserService{
        users:  make(map[int]*User),
        nextID: 1,
    }
}

func (us *UserService) Create(user *User) *User {
    user.ID = us.nextID
    us.nextID++
    user.Created = time.Now()
    user.Modified = time.Now()
    us.users[user.ID] = user
    return user
}

func (us *UserService) GetByID(id int) (*User, error) {
    if user, exists := us.users[id]; exists {
        return user, nil
    }
    return nil, echo.NewHTTPError(http.StatusNotFound, "User not found")
}

func (us *UserService) Update(id int, updates *User) (*User, error) {
    user, err := us.GetByID(id)
    if err != nil {
        return nil, err
    }

    if updates.Name != "" {
        user.Name = updates.Name
    }
    if updates.Email != "" {
        user.Email = updates.Email
    }
    user.Modified = time.Now()

    return user, nil
}

func (us *UserService) Delete(id int) (*User, error) {
    user, err := us.GetByID(id)
    if err != nil {
        return nil, err
    }

    delete(us.users, id)
    return user, nil
}

func (us *UserService) List() []*User {
    users := make([]*User, 0, len(us.users))
    for _, user := range us.users {
        users = append(users, user)
    }
    return users
}

// Handlers
func createUserHandler(userService *UserService) echo.HandlerFunc {
    return func(c echo.Context) error {
        ctx := c.Request().Context()

        var user User
        if err := c.Bind(&user); err != nil {
            return err
        }

        created := userService.Create(&user)
        
        // Record creation in audit context
        core.AuditableCreate(ctx, created)

        return c.JSON(http.StatusCreated, created)
    }
}

func getUserHandler(userService *UserService) echo.HandlerFunc {
    return func(c echo.Context) error {
        ctx := c.Request().Context()
        
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
        }

        user, err := userService.GetByID(id)
        if err != nil {
            return err
        }

        // Record read operation
        core.AuditableGet(ctx, id)

        return c.JSON(http.StatusOK, user)
    }
}

func updateUserHandler(userService *UserService) echo.HandlerFunc {
    return func(c echo.Context) error {
        ctx := c.Request().Context()
        
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
        }

        // Get existing user for audit trail
        existingUser, err := userService.GetByID(id)
        if err != nil {
            return err
        }

        var updates User
        if err := c.Bind(&updates); err != nil {
            return err
        }

        updatedUser, err := userService.Update(id, &updates)
        if err != nil {
            return err
        }

        // Record update operation with before/after data
        core.AuditableUpdate(ctx, id, existingUser, updatedUser)

        return c.JSON(http.StatusOK, updatedUser)
    }
}

func deleteUserHandler(userService *UserService) echo.HandlerFunc {
    return func(c echo.Context) error {
        ctx := c.Request().Context()
        
        id, err := strconv.Atoi(c.Param("id"))
        if err != nil {
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
        }

        deletedUser, err := userService.Delete(id)
        if err != nil {
            return err
        }

        // Record deletion with deleted data
        core.AuditableDelete(ctx, id, deletedUser)

        return c.JSON(http.StatusOK, map[string]string{
            "message": "User deleted successfully",
        })
    }
}

func listUsersHandler(userService *UserService) echo.HandlerFunc {
    return func(c echo.Context) error {
        ctx := c.Request().Context()
        
        users := userService.List()

        // Record list operation with metadata
        core.AuditableList(ctx, map[string]interface{}{
            "total_count": len(users),
            "query_time": time.Now(),
        })

        return c.JSON(http.StatusOK, users)
    }
}

// Extract user ID from context (for demo purposes)
func extractUserFromContext(ctx context.Context) any {
    // In a real application, this would extract from JWT, session, etc.
    if userID := ctx.Value("user_id"); userID != nil {
        return userID
    }
    return "anonymous" // Default user
}

// Skip auditing for health checks
func shouldSkipAudit(ctx context.Context) bool {
    // In a real application, check request path, headers, etc.
    return false // Audit everything for demo
}

func main() {
    // Initialize services
    userService := NewUserService()
    
    // Create audit configuration
    auditConfig := transport.Config{
        Auditor: &ConsoleAuditor{
            logger: log.New(log.Writer(), "AUDIT ", log.LstdFlags),
        },
        NewTraceID: func() any {
            return uuid.New().String()
        },
        UserFromContext: extractUserFromContext,
        Skipper:        shouldSkipAudit,
    }

    // Initialize Echo with audit middleware
    e := echo.New()
    
    // Basic middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    
    // Configure audit system
    echoadapter.Configure(auditConfig)
    e.Use(echoadapter.Root())

    // Setup routes with audit configuration
    userAudit := echoadapter.For(EntityUser)
    
    userGroup := e.Group("/users")
    userGroup.POST("", createUserHandler(userService), userAudit(ActionCreate))
    userGroup.GET("/:id", getUserHandler(userService), userAudit(ActionRead))
    userGroup.PUT("/:id", updateUserHandler(userService), userAudit(ActionUpdate))
    userGroup.DELETE("/:id", deleteUserHandler(userService), userAudit(ActionDelete))
    userGroup.GET("", listUsersHandler(userService), userAudit(ActionList))

    // Health check endpoint (could be configured to skip auditing)
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{
            "status": "healthy",
            "time":   time.Now().Format(time.RFC3339),
        })
    })

    // Start server
    log.Println("🚀 Starting server on :8080")
    log.Println("📝 Try these endpoints:")
    log.Println("  POST   /users")
    log.Println("  GET    /users/:id")
    log.Println("  PUT    /users/:id")
    log.Println("  DELETE /users/:id")
    log.Println("  GET    /users")
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

**Test the application:**

```bash
# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Get user
curl http://localhost:8080/users/1

# Update user
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com"}'

# List users
curl http://localhost:8080/users

# Delete user
curl -X DELETE http://localhost:8080/users/1
```

---

## Best Practices

### 1. **Error Handling**
```go
// Always handle audit errors gracefully
func (em *EchoMw) publish(ctx context.Context, auditCtx *core.AuditableContext) {
    event := auditCtx.ToEvent()
    
    if err := em.cfg.Auditor.Audit(ctx, *event); err != nil {
        // Log error but don't fail the request
        log.Printf("Failed to publish audit event: %v", err)
        
        // Optional: Send to dead letter queue or retry mechanism
    }
}
```

### 2. **Sensitive Data Protection**
```go
// Sanitize sensitive data before auditing
func sanitizeUser(user *User) *User {
    sanitized := *user
    sanitized.Password = "[REDACTED]"
    sanitized.SSN = "[REDACTED]"
    return &sanitized
}

// Use in handlers
core.AuditableCreate(ctx, sanitizeUser(user))
```

### 3. **Custom Metadata**
```go
// Add contextual information to events
core.SetContext(ctx,
    core.WithResourceID(userID),
    core.WithMetadata("department", "engineering"),
    core.WithMetadata("api_version", "v2"),
    core.WithMetadata("feature_flag", "new_user_flow"),
)
```

---

## License

This project is licensed under the [MIT License](LICENSE).

---

**Note:** This documentation reflects the actual structure and functionality of your audit library. All code examples are functional and ready to use. For specific implementation details, refer to the source code in the respective packages.