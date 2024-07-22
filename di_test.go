package godi

import (
	"testing"
)

type TestService struct {
	Value string
}

type DependentService struct {
	TestService *TestService
}

func TestRegisterAndResolve(t *testing.T) {
	container := New()

	Register(container, func() *TestService {
		return &TestService{Value: "test"}
	}, Transient)

	service, err := Resolve[*TestService](container)
	if err != nil {
		t.Fatalf("Failed to resolve TestService: %v", err)
	}

	if service.Value != "test" {
		t.Errorf("Expected Value to be 'test', got '%s'", service.Value)
	}
}

func TestSingletonLifecycle(t *testing.T) {
	container := New()

	Register(container, func() *TestService {
		return &TestService{Value: "singleton"}
	}, Singleton)

	service1, _ := Resolve[*TestService](container)
	service2, _ := Resolve[*TestService](container)

	if service1 != service2 {
		t.Error("Expected singleton services to be the same instance")
	}
}

func TestTransientLifecycle(t *testing.T) {
	container := New()

	Register(container, func() *TestService {
		return &TestService{Value: "transient"}
	}, Transient)

	service1, _ := Resolve[*TestService](container)
	service2, _ := Resolve[*TestService](container)

	if service1 == service2 {
		t.Error("Expected transient services to be different instances")
	}
}

func TestDependencyResolution(t *testing.T) {
	container := New()

	Register(container, func() *TestService {
		return &TestService{Value: "dependency"}
	}, Singleton)

	Register(container, func() *DependentService {
		testService, _ := Resolve[*TestService](container)
		return &DependentService{TestService: testService}
	}, Transient)

	dependent, err := Resolve[*DependentService](container)
	if err != nil {
		t.Fatalf("Failed to resolve DependentService: %v", err)
	}

	if dependent.TestService == nil {
		t.Error("Expected DependentService to have a non-nil TestService")
	}

	if dependent.TestService.Value != "dependency" {
		t.Errorf("Expected TestService Value to be 'dependency', got '%s'", dependent.TestService.Value)
	}
}

func TestMustResolve(t *testing.T) {
	container := New()

	Register(container, func() *TestService {
		return &TestService{Value: "must resolve"}
	}, Transient)

	service := MustResolve[*TestService](container)
	if service.Value != "must resolve" {
		t.Errorf("Expected Value to be 'must resolve', got '%s'", service.Value)
	}
}

func TestMustResolvePanic(t *testing.T) {
	container := New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected MustResolve to panic for unregistered service")
		}
	}()

	MustResolve[*TestService](container) // This should panic
}