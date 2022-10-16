package genjector

import (
	"fmt"
)

type Key struct {
	Annotation string
	Value      interface{}
}

func (k Key) Generate() interface{} {
	if len(k.Annotation) > 0 {
		return [2]interface{}{k.Annotation, k.Value}
	}
	return k.Value
}

type KeySource interface {
	Key() Key
}

type KeyOption interface {
	Key(key Key) Key
	Container(container Container) Container
}

type Binding interface {
	Instance(initialize bool) (interface{}, error)
}

type BindingSource[T any] interface {
	Key() Key
	Binding() (Binding, error)
}

type FollowingBindingSource[T any] interface {
	SetPrevious(binding Binding)
}

type BindingOption interface {
	Key(key Key) Key
	Container(container Container) Container
	Binding(binding Binding) (Binding, error)
}

type Container map[interface{}]Binding

var global = NewContainer()

func NewContainer() Container {
	return map[interface{}]Binding{}
}

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

func MustBind[T any](source BindingSource[T], options ...BindingOption) {
	err := Bind(source, options...)
	if err != nil {
		panic(err)
	}
}

func Initialize[T any](options ...KeyOption) (T, error) {
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

func MustInitialize[T any](options ...KeyOption) T {
	instance, err := Initialize[T](options...)
	if err != nil {
		panic(err)
	}

	return instance
}

func Clean() {
	global = NewContainer()
}

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
