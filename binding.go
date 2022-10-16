package genjector

import "fmt"

// bindingSource is a concrete implementation for BindingSource interface.
type bindingSource[T any] struct {
	binding   Binding
	keySource baseKeySource[T]
}

// Binding returns containing instance of Binding interface. Initially it makes
// the concrete instance, to check if instance matches desired type of Binding.
//
// It respects BindingSource interface.
func (s *bindingSource[T]) Binding() (Binding, error) {
	instance, err := s.binding.Instance(false)
	if err != nil {
		return nil, err
	}

	if _, ok := instance.(T); !ok {
		var initial T
		return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, initial, instance)
	}
	return s.binding, nil
}

// Key executes the same method from inner KeyOption instance.
//
// It respects BindingSource interface.
func (s *bindingSource[T]) Key() Key {
	return s.keySource.Key()
}

// Initializable represents any struct that contains a method Init.
// When such struct as defined AsReference or AsValue, method Init will be
// called during initialization process.
type Initializable interface {
	Init()
}

// valueBinding is a concrete implementation for Binding interface.
type valueBinding[S any] struct{}

// Instance delivers the value of the concrete instance of type S.
// If the reference to the struct respects Initializable interface,
// Init method will be called.
//
// It respects Binding interface.
func (valueBinding[S]) Instance(initialize bool) (interface{}, error) {
	initial := *new(S)
	var instance interface{} = &initial
	if !initialize {
		return initial, nil
	}

	if value, ok := instance.(Initializable); ok {
		value.Init()
	}

	return initial, nil
}

// AsValue delivers a BindingSource for a type T, by binding a value of a struct
// to the concrete interface (or the struct itself). It must be only used with value and
// not reference. In case reference is used, code will return a nil value for the instance.
//
// Example:
// err := genjector.Bind(genjector.AsValue[ValueInterface, ValueStruct]())
//
// BindingSource can be only used as the first argument to Bind method.
func AsValue[T any, S any]() BindingSource[T] {
	return &bindingSource[T]{
		binding:   valueBinding[S]{},
		keySource: baseKeySource[T]{},
	}
}

// referenceBinding is a concrete implementation for Binding interface.
type referenceBinding[R any] struct{}

// Instance delivers the reference of the concrete instance of type S.
// If the struct respects Initializable interface, Init method will be called.
//
// It respects Binding interface.
func (referenceBinding[R]) Instance(initialize bool) (interface{}, error) {
	var instance interface{} = new(R)
	if !initialize {
		return instance, nil
	}

	if value, ok := instance.(Initializable); ok {
		value.Init()
	}

	return instance, nil
}

// AsReference delivers a BindingSource for a type T, by binding reference of a struct
// to the concrete interface (or the struct itself). It must be only used with references and
// not values. In case values is used, code will panic.
//
// Example:
// err := genjector.Bind(genjector.AsReference[ReferenceInterface, *ReferenceStruct]())
//
// BindingSource can be only used as the first argument to Bind method.
func AsReference[T any, S *R, R any]() BindingSource[T] {
	return &bindingSource[T]{
		binding:   referenceBinding[R]{},
		keySource: baseKeySource[T]{},
	}
}

// ProviderMethod defines a type of a method that should delivers
// an instance od type S. This method acts as an constructor method
// and it is executed at the time of Initialize method.
//
// It respects Binding interface.
type ProviderMethod[S any] func() (S, error)

// Instance delivers the concrete instance of type S, by executing
// root ProviderMethod itself.
//
// It respects Binding interface.
func (s ProviderMethod[S]) Instance(bool) (interface{}, error) {
	return s()
}

// AsProvider delivers a BindingSource for a type T, by defining a ProviderMethod
// (or constructor method) for the new instance of some interface (or a struct).
//
// Example:
//
//	err := genjector.Bind(genjector.AsProvider[ProviderInterface](func() (*ProviderStruct, error) {
//	  return &ProviderStruct{
//	    value: "value provided inside the ProviderMethod",
//	  }, nil
//	}))
//
// BindingSource can be only used as the first argument to Bind method.
func AsProvider[T any, S any](provider ProviderMethod[S]) BindingSource[T] {
	return &bindingSource[T]{
		binding:   provider,
		keySource: baseKeySource[T]{},
	}
}

// instanceBinding is a concrete implementation for Binding interface.
type instanceBinding[S any] struct {
	instance S
}

// Instance delivers the concrete instance of type S, by returning already
// initialized instance that instanceBinding holds.
//
// It respects Binding interface.
func (s *instanceBinding[S]) Instance(bool) (interface{}, error) {
	return s.instance, nil
}

// AsInstance delivers a BindingSource for a type T, by using a concrete
// instance that is passed as an argument to AsInstance method, to returns
// that instance whenever it is required from Binding.
//
// Example:
//
//	err := genjector.Bind(genjector.AsInstance[*InstanceStruct](&InstanceStruct{
//	  value: "value provided in concrete instance",
//	}))
//
// BindingSource can be only used as the first argument to Bind method.
func AsInstance[T any, S any](instance S) BindingSource[T] {
	return &bindingSource[T]{
		binding: &instanceBinding[S]{
			instance: instance,
		},
		keySource: baseKeySource[T]{},
	}
}

// bindingOption is a concrete implementation for BindingOption interface.
type bindingOption struct {
	bindingFunc func(binding Binding) (Binding, error)
	keyOption   KeyOption
}

// Binding executes the inner bindingFunc method.
//
// It respects BindingOption interface.
func (b *bindingOption) Binding(binding Binding) (Binding, error) {
	return b.bindingFunc(binding)
}

// Key executes the same method from inner KeyOption instance.
//
// It respects BindingOption interface.
func (b *bindingOption) Key(key Key) Key {
	return b.keyOption.Key(key)
}

// Container executes the same method from inner KeyOption instance.
//
// It respects BindingOption interface.
func (b *bindingOption) Container(container Container) Container {
	return b.keyOption.Container(container)
}

// singletonBinding is a concrete implementation for Binding interface.
type singletonBinding struct {
	parent      Binding
	singleton   interface{}
	initialized bool
}

// Instance delivers already stored instance, which should be present if this
// method was already executed before. Otherwise it retrieves the instance from
// a child Binding and stores it internally for the next calls.
//
// It respects Binding interface.
func (b *singletonBinding) Instance(initialize bool) (interface{}, error) {
	if b.initialized {
		return b.singleton, nil
	}

	instance, err := b.parent.Instance(initialize)
	if err != nil {
		return nil, err
	}

	b.initialized = true
	b.singleton = instance
	return instance, nil
}

// AsSingleton delivers a BindingOption that defines the instance of desired
// Binding as a singleton. That means only first time the Init method (or ProviderMethod)
// will be called, and every next time the same instance will be delivered
// as a result of Initialize method.
//
// Example:
// err := genjector.Bind(
//
//	genjector.AsReference[SingletonInterface, *SingletonStruct](),
//	genjector.AsSingleton(),
//
// )
//
// AsSingleton should be only used as a BindingOption for Bind method, as it
// does not affect functionality if it is used in Initialize method.
func AsSingleton() BindingOption {
	return &bindingOption{
		bindingFunc: func(binding Binding) (Binding, error) {
			return &singletonBinding{
				parent: binding,
			}, nil
		},
		keyOption: sameKeyOption{},
	}
}

// WithAnnotation delivers a BindingOption that allows to name specific Binding
// with any annotation desired.
//
// Example:
// err = genjector.Bind(
//
//	genjector.AsReference[AnnotationInterface, *AnnotationStruct](),
//	genjector.WithAnnotation("first"),
//
// )
//
// To properly use a customer Container, WithAnnotation should be used in both
// Bind and Initialize methods.
func WithAnnotation(annotation string) BindingOption {
	return &bindingOption{
		bindingFunc: func(binding Binding) (Binding, error) {
			return binding, nil
		},
		keyOption: &annotatedKeyOption{
			annotation: annotation,
		},
	}
}

// WithContainer delivers a BindingOption that overrides the usage of standard
// internal (global) Container. It allows to provide a fresh, a custom instance
// of Container, that can be made from NewContainer method.
//
// Example:
// err := genjector.Bind(
//
//	genjector.AsReference[ContainerInterface, *ContainerStruct](),
//	genjector.WithContainer(customContainer),
//
// )
//
// To properly use a customer Container, WithContainer should be used in both
// Bind and Initialize methods.
func WithContainer(container Container) BindingOption {
	return &bindingOption{
		bindingFunc: func(binding Binding) (Binding, error) {
			return binding, nil
		},
		keyOption: &containerKeyOption{
			container: container,
		},
	}
}
