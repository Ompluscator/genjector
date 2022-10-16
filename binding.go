package genjector

import "fmt"

type bindingSource[T any] struct {
	binding   Binding
	keySource baseKeySource[T]
}

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

func (s *bindingSource[T]) Key() Key {
	return s.keySource.Key()
}

type valueBinding[S any] struct{}

func (valueBinding[S]) Instance(initialize bool) (interface{}, error) {
	initial := *new(S)
	var instance interface{} = &initial
	if !initialize {
		return initial, nil
	}

	if value, ok := instance.(initializable); ok {
		value.Init()
	}

	return initial, nil
}

func AsValue[T any, S any]() BindingSource[T] {
	return &bindingSource[T]{
		binding:   valueBinding[S]{},
		keySource: baseKeySource[T]{},
	}
}

type initializable interface {
	Init()
}

type referenceBinding[R any] struct{}

func (referenceBinding[R]) Instance(initialize bool) (interface{}, error) {
	var instance interface{} = new(R)
	if !initialize {
		return instance, nil
	}

	if value, ok := instance.(initializable); ok {
		value.Init()
	}

	return instance, nil
}

func AsReference[T any, S *R, R any]() BindingSource[T] {
	return &bindingSource[T]{
		binding:   referenceBinding[R]{},
		keySource: baseKeySource[T]{},
	}
}

type ProviderMethod[S any] func() (S, error)

func (b ProviderMethod[S]) Instance(bool) (interface{}, error) {
	return b()
}

func AsProvider[T any, S any](provider ProviderMethod[S]) BindingSource[T] {
	return &bindingSource[T]{
		binding:   provider,
		keySource: baseKeySource[T]{},
	}
}

type instanceBinding[S any] struct {
	instance S
}

func (b *instanceBinding[S]) Instance(bool) (interface{}, error) {
	return b.instance, nil
}

func AsInstance[T any, S any](instance S) BindingSource[T] {
	return &bindingSource[T]{
		binding: &instanceBinding[S]{
			instance: instance,
		},
		keySource: baseKeySource[T]{},
	}
}

type bindingOption struct {
	bindingFunc func(binding Binding) (Binding, error)
	keyOption   KeyOption
}

func (b *bindingOption) Binding(binding Binding) (Binding, error) {
	return b.bindingFunc(binding)
}

func (b *bindingOption) Key(key Key) Key {
	return b.keyOption.Key(key)
}

func (b *bindingOption) Container(container Container) Container {
	return b.keyOption.Container(container)
}

type singletonBinding struct {
	parent      Binding
	singleton   interface{}
	initialized bool
}

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
