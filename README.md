# 🧙‍♂️ Godi - Go Dependency Injection Engine

A lightweight, type-safe Dependency Injection (DI) container for Go, leveraging generics for a clean and intuitive API.

## ✨ Features

- Type-safe dependency registration and resolution using Go generics
- Support for both singleton and transient lifecycles
- Thread-safe operations
- No reflection for improved performance
- Simple and intuitive API

## 🚀 Installation

To install the Godi Engine, use `go get`:

```bash
go get github.com/devs-group/godi
```

## 🏁 Quick Start

Here's a simple example of how to use the Godi Engine:

```go
package main

import (
	"fmt"
	"github.com/devs-group/godi"
)

type Database struct {
	ConnectionString string
}

type UserRepository struct {
	DB *Database
}

type UserService struct {
	Repo *UserRepository
}

func main() {
	container := godi.New()

	// Register services
	godi.Register(container, func() *Database {
		return &Database{ConnectionString: "example_connection_string"}
	}, godi.Singleton)

	godi.Register(container, func() *UserRepository {
		db, _ := godi.Resolve[*Database](container)
		return &UserRepository{DB: db}
	}, godi.Transient)

	godi.Register(container, func() *UserService {
		repo, _ := godi.Resolve[*UserRepository](container)
		return &UserService{Repo: repo}
	}, godi.Transient)

	// Resolve and use a service
	userService, _ := godi.Resolve[*UserService](container)
	fmt.Printf("Connection string: %s\n", userService.Repo.DB.ConnectionString)
}
```

## 📚 API Reference

### Creating a Container

```go
container := godi.New()
```

### Registering Services

```go
godi.Register[T any](container *Container, constructor func() T, lifecycle Lifecycle)
```

- `container`: The DI container
- `constructor`: A function that creates an instance of the service
- `lifecycle`: Either `godi.Singleton` or `godi.Transient`

### Resolving Services

```go
service, err := godi.Resolve[T any](container *Container)
```

### Must Resolve (Panics on Error)

```go
service := godi.MustResolve[T any](container *Container)
```

## 🔄 Lifecycle Management

- **Singleton**: Only one instance is created and reused for all subsequent resolves.
- **Transient**: A new instance is created each time the service is resolved.

## 🔒 Thread Safety

The Godi Engine is designed to be thread-safe. You can use a single container across multiple goroutines without worrying about race conditions.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.