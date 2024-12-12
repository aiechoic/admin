package ioc

import (
	"context"
	"io"
	"sync"
)

// injector is an interface for creating new instances.
type injector interface {
	new(c *Container) (any, error)
}

// Provider is a provider for a type T.
type Provider[T any] struct {
	fn func(c *Container) (T, error)
	mu sync.Mutex
}

func (p *Provider[T]) new(c *Container) (any, error) {
	return p.fn(c)
}

// NewProvider creates a new provider for the given type T.
func NewProvider[T any](fn func(c *Container) (T, error)) *Provider[T] {
	return &Provider[T]{fn: fn}
}

// Get returns an instance of T. If the instance is already created, it returns the existing instance.
// If the instance is not created, it creates a new instance, stores it, and returns it.
func (p *Provider[T]) Get(c *Container) (t T, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, ok := c.get(p)
	if ok {
		return v.(T), nil
	}
	t, err = p.fn(c)
	if err != nil {
		return
	}
	c.set(p, t)
	return t, nil
}

// MustGet is like Get, but it panics if an error occurs.
func (p *Provider[T]) MustGet(c *Container) T {
	t, err := p.Get(c)
	if err != nil {
		panic(err)
	}
	return t
}

// GetNew creates a new instance of T and returns it. It does not store the instance.
func (p *Provider[T]) GetNew(c *Container) (T, error) {
	t, err := p.fn(c)
	if err != nil {
		var zero T
		return zero, err
	}
	return t, nil
}

// MustGetNew is like GetNew, but it panics if an error occurs.
func (p *Provider[T]) MustGetNew(c *Container) T {
	t, err := p.GetNew(c)
	if err != nil {
		panic(err)
	}
	return t
}

// IsSet returns true if the instance of T is already created.
func (p *Provider[T]) IsSet(c *Container) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	_, ok := c.get(p)
	return ok
}

// Set stores the instance of T.
func (p *Provider[T]) Set(c *Container, v T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	c.set(p, v)
}

type Generator[T any, C comparable] func(name string, args ...C) *Provider[T]

// Providers is a collection of providers. It is used to manage named providers.
type Providers[T any, C comparable] struct {
	ps   map[string]*Provider[T]
	args map[string][]C
	fn   Generator[T, C]
	mu   sync.Mutex
}

// NewProviders creates a new Providers instance.
// The fn parameter is a function that creates a new provider for the given name.
func NewProviders[T any, C comparable](fn Generator[T, C]) *Providers[T, C] {
	return &Providers[T, C]{
		ps:   map[string]*Provider[T]{},
		args: map[string][]C{},
		fn:   fn,
	}
}

// GetProvider returns a provider with the given name.
func (ps *Providers[T, C]) GetProvider(name string, args ...C) *Provider[T] {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	p, ok := ps.ps[name]
	if !ok {
		p = ps.fn(name, args...)
		ps.ps[name] = p
		ps.args[name] = args
	} else {
		_args := ps.args[name]
		if len(args) != len(_args) {
			panic("args changed")
		}
		for i, dep := range args {
			if dep != _args[i] {
				panic("args changed")
			}
		}
	}
	return p
}

// Container is a container for managing instances.
type Container struct {
	instances map[injector]any
	cancel    context.CancelFunc
	mu        sync.Mutex
}

// NewContainer creates a new container.
func NewContainer() *Container {
	return &Container{
		instances: map[injector]any{},
	}
}

func (c *Container) get(p injector) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	i, ok := c.instances[p]
	return i, ok
}

func (c *Container) set(p injector, ins any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.instances[p] = ins
}

// GetAllInstances returns all instances in the container.
func (c *Container) GetAllInstances() []any {
	c.mu.Lock()
	defer c.mu.Unlock()

	var instances []any
	for _, ins := range c.instances {
		instances = append(instances, ins)
	}
	return instances
}

// Errors is a collection of errors.
type Errors []error

func (errs Errors) Error() string {
	var str string
	for _, err := range errs {
		str += err.Error() + "\n"
	}
	return str
}

// Close closes all instances that implement the io.Closer interface.
func (c *Container) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	var errs Errors
	for _, ins := range c.instances {
		if closer, ok := ins.(io.Closer); ok {
			err := closer.Close()
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
