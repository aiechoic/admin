package ioc

import (
	"errors"
	"fmt"
)

// ServiceA is a simple service that depends on ServiceB
type ServiceA struct {
	ServiceB *ServiceB
}

// ServiceB is another simple service
type ServiceB struct{}

func (s *ServiceB) Close() error {
	return errors.New("ServiceB.Close error")
}

var ServiceBProvider = NewProvider(func(c *Container) (*ServiceB, error) {
	return &ServiceB{}, nil
})

var ServiceAProvider = NewProvider(func(c *Container) (*ServiceA, error) {
	serviceB, err := ServiceBProvider.Get(c)
	if err != nil {
		return nil, err
	}
	return &ServiceA{ServiceB: serviceB}, nil
})

func ExampleProvider_Get() {
	c := NewContainer()

	defer func() {
		err := c.Close()
		if err != nil {
			fmt.Println(err) // Output: ServiceB.Close error
		}
	}()

	serviceA, err := ServiceAProvider.Get(c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("serviceA.ServiceB == nil: %v\n", serviceA.ServiceB == nil)

	serviceB, err := ServiceBProvider.Get(c)
	if err != nil {
		panic(err)
	}
	fmt.Printf("serviceA.ServiceB == serviceB: %v\n", serviceA.ServiceB == serviceB)

	// Output:
	// serviceA.ServiceB == nil: false
	// serviceA.ServiceB == serviceB: true
	// ServiceB.Close error
}

type Client struct {
	Name string
}

var clientProviders = NewProviders(func(name string, args ...any) *Provider[*Client] {
	return NewProvider(func(c *Container) (*Client, error) {
		return &Client{Name: name}, nil
	})
})

func ExampleNewProviders() {
	c := NewContainer()

	client1 := clientProviders.GetProvider("client1").MustGet(c)
	client2 := clientProviders.GetProvider("client2").MustGet(c)
	clientOld := clientProviders.GetProvider("client1").MustGet(c)

	fmt.Printf("client1.Name='%s'\n", client1.Name)
	fmt.Printf("client2.Name='%s'\n", client2.Name)
	fmt.Printf("client1 == clientOld: %v\n", client1 == clientOld)

	// Output:
	// client1.Name='client1'
	// client2.Name='client2'
	// client1 == clientOld: true
}
