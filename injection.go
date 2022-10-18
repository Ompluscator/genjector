package genjector

import (
	"fmt"
)

// Key is a struct that contains information for Binding keys
// inside a Container.
//
// It is meant to be used only for internal purposes.
type Key struct {
	Annotation string
	Value      interface{}
}

// Generate delivers a final Binding key for the Container.
func (k Key) Generate() interface{} {
	if len(k.Annotation) > 0 {
		return [2]interface{}{k.Annotation, k.Value}
	}
	return k.Value
}

// KeySource represents an interface that builds a Key for Binding.
type KeySource interface {
	Key() Key
}

// KeyOption represents an interface that overrides creation of Key
// and Container.
type KeyOption interface {
	Key(key Key) Key
	Container(container Container) Container
}

// Binding represents an interface that delivers new instance for
// particular interface (or a struct).
type Binding interface {
	Instance(initialize bool) (interface{}, error)
}

// BindingSource represents an interface that delivers starting Key and
// Binding instances, that later could be overridden by KeyOption or
// BindingOption.
type BindingSource[T any] interface {
	Key() Key
	Binding() (Binding, error)
}

// FollowingBindingSource represents an interface for a Binding instance
// that requires to get previous instance of Binding inside Container,
// before the new one should be stored on that place.
type FollowingBindingSource[T any] interface {
	SetPrevious(binding Binding)
}

// BindingOption represents an interface that overrides creation of Key,
// Binding and Container.
type BindingOption interface {
	Key(key Key) Key
	Container(container Container) Container
	Binding(binding Binding) (Binding, error)
}

// Container is a child type used for storing all Binding instances.
type Container map[interface{}]Binding

// global is a concrete global Container
var global = NewContainer()

// NewContainer delivers a new instance of Container.
func NewContainer() Container {
	return map[interface{}]Binding{}
}

// Bind executes complete logic for binding particular value (or pointer) to
// desired interface (or struct). By default, it stores all Binding instances
// into default inner Container.
//
// It requires only BindingSource to be passed as an argument and all other
// instances of BindingOption are optional.
//
// At this point, if Binding for particular interface (or struct) is not defined,
// it uses its own fallback Binding. Still, it works fully only for values, not pointers,
// as for pointers it returns nil value. That means that pointer Binding
// should be always defined.
func Bind[T any](source BindingSource[T], options ...BindingOption) error {
	key := source.Key()

	internal := global
	for _, option := range options {
		key = option.Key(key)
		internal = option.Container(internal)
	}

	generated := key.Generate()

	if child, ok := source.(FollowingBindingSource[T]); ok {
		parent, ok := internal[generated]
		if ok {
			child.SetPrevious(parent)
		}
	}

	binding, err := source.Binding()
	if err != nil {
		return err
	}

	for _, option := range options {
		binding, err = option.Binding(binding)
		if err != nil {
			return err
		}
	}

	internal[generated] = binding
	return nil
}

// MustBind wraps Bind method, by making sure error is not returned as an argument.
//
// Still, in case of error, it panics.
func MustBind[T any](source BindingSource[T], options ...BindingOption) {
	err := Bind(source, options...)
	if err != nil {
		panic(err)
	}
}

// NewInstance executes complete logic for initializing value (or pointer) for
// desired interface (or struct). By default, it uses Binding instance from default
// inner Container. If such Binding can not be found, it tries to make its own
// fallback Binding.
//
// All instances of BindingOption are optional.
//
// At this point, if Binding for particular interface (or struct) is not defined,
// it uses its own fallback Binding. Still, it works fully only for values, not pointers,
// as for pointers it returns nil value. That means that pointer Binding
// should be always defined.
func NewInstance[T any](options ...KeyOption) (T, error) {
	var empty T
	source := &baseKeySource[T]{}

	key := source.Key()

	internal := global
	for _, option := range options {
		key = option.Key(key)
		internal = option.Container(internal)
	}

	generated := key.Generate()

	binding, ok := internal[generated]
	if !ok {
		var err error
		binding, err = getFallbackBinding[T]()
		if err != nil {
			return empty, err
		}
	}

	instance, err := binding.Instance(true)
	if err != nil {
		return empty, err
	}

	result, ok := instance.(T)
	if !ok {
		return empty, fmt.Errorf(`invalid binding is defined for key "%v"`, generated)
	}

	return result, nil
}

// MustNewInstance wraps NewInstance method, by making sure error is not returned as an argument.
//
// Still, in case of error, it panics.
func MustNewInstance[T any](options ...KeyOption) T {
	instance, err := NewInstance[T](options...)
	if err != nil {
		panic(err)
	}

	return instance
}

// Clean creates a new instance of inner Container.
func Clean() {
	global = NewContainer()
}

// getFallbackBinding creates a new instance of fallback Binding.
func getFallbackBinding[T any]() (Binding, error) {
	var binding Binding
	var err error

	source := AsValue[T, T]()
	binding, err = source.Binding()
	if err == nil {
		instance, err := binding.Instance(false)
		if err == nil {
			if _, ok := instance.(T); ok {
				return binding, nil
			}
		}
	}

	return nil, err
}
