---
applyTo: ["*"]
---

# Service Registration Rules

## Registration Methods

- **Singleton**: `ctn.RegisterSingleton("name", func() interface{} { return instance })` - For stateless services (database, config, cache)
- **Transient**: `ctn.RegisterTransient("name", func() interface{} { return instance })` - For stateful services (business logic, handlers)
- **Instance**: `ctn.Register("name", instance)` - For pre-created objects

## Registration Order

1. Core services (config, logger, container)
2. Infrastructure (database, cache, external APIs)
3. Repositories (data access layer)
4. Services (business logic layer)
5. Handlers (HTTP handlers)
6. Make sure to validate services dependencies using container/validation.go

## Service Retrieval

```go
// Type-safe (recommended)
var service *ServiceType
err := container.GetTyped("serviceName", &service)

// Manual
service, err := container.Get("serviceName")
serviceTyped := service.(*ServiceType)
```

## Lifecycle Guidelines

- **Singleton**: Database connections, configuration, cache clients
- **Transient**: Business services, HTTP handlers, processors
- **Instance**: Pre-created objects, third-party clients

## Dependencies

```go
ctn.RegisterTransient("userService", func() interface{} {
    repo, _ := ctn.Get("userRepo")
    logger, _ := ctn.Get("logger")
    return &UserService{
        repo: repo.(*UserRepository),
        logger: logger.(*logger.Logger),
    }
})
```
