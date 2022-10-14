package genjector

import "fmt"

type sliceBinding[T any] struct {
	previous *sliceBinding[T]
	current  Binding
}

func (b *sliceBinding[T]) Instance(initialize bool) (interface{}, error) {
	var result []T
	if initialize && b.previous != nil {
		instance, err := b.previous.Instance(initialize)
		if err != nil {
			return nil, err
		}

		transformed, ok := instance.([]T)
		if !ok {
			return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, result, instance)
		}
		result = append(result, transformed...)
	}

	instance, err := b.current.Instance(initialize)
	if err != nil {
		return nil, err
	}

	transformed, ok := instance.(T)
	if !ok {
		return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, result, instance)
	}

	result = append(result, transformed)
	return result, err
}

type sliceBindingSource[T any] struct {
	parent    Binding
	source    BindingSource[T]
	keySource KeySource
}

func (s *sliceBindingSource[T]) Binding() (Binding, error) {
	binding, err := s.source.Binding()
	if err != nil {
		return nil, err
	}

	instance, _ := binding.Instance(false)
	if _, ok := instance.(T); !ok {
		var initial T
		return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, initial, instance)
	}

	previous, ok := s.parent.(*sliceBinding[T])
	if !ok {
		return &sliceBinding[T]{
			current: binding,
		}, nil
	}

	return &sliceBinding[T]{
		previous: previous,
		current:  binding,
	}, nil
}

func (s *sliceBindingSource[T]) SetPrevious(binding Binding) {
	s.parent = binding
}

func (s *sliceBindingSource[T]) Key() Key {
	return s.keySource.Key()
}

func InSlice[T any](source BindingSource[T]) BindingSource[T] {
	return &sliceBindingSource[T]{
		source:    source,
		keySource: &sliceKeySource[T]{},
	}
}

type mapBinding[K comparable, T any] struct {
	previous *mapBinding[K, T]
	key      K
	current  Binding
}

func (b *mapBinding[K, T]) Instance(initialize bool) (interface{}, error) {
	result := map[K]T{}
	if initialize && b.previous != nil {
		instance, err := b.previous.Instance(initialize)
		if err != nil {
			return nil, err
		}

		transformed, ok := instance.(map[K]T)
		if !ok {
			return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, result, instance)
		}
		result = transformed
	}

	instance, err := b.current.Instance(initialize)
	if err != nil {
		return nil, err
	}

	transformed, ok := instance.(T)
	if !ok {
		return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, result, instance)
	}

	result[b.key] = transformed
	return result, err
}

type mapBindingSource[K comparable, T any] struct {
	parent    Binding
	source    BindingSource[T]
	key       K
	keySource KeySource
}

func (s *mapBindingSource[K, T]) Binding() (Binding, error) {
	binding, err := s.source.Binding()
	if err != nil {
		return nil, err
	}

	instance, _ := binding.Instance(false)
	if _, ok := instance.(T); !ok {
		var initial T
		return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, initial, instance)
	}

	previous, ok := s.parent.(*mapBinding[K, T])
	if !ok {
		return &mapBinding[K, T]{
			key:     s.key,
			current: binding,
		}, nil
	}

	return &mapBinding[K, T]{
		previous: previous,
		key:      s.key,
		current:  binding,
	}, nil
}

func (s *mapBindingSource[K, T]) SetPrevious(binding Binding) {
	s.parent = binding
}

func (s *mapBindingSource[K, T]) Key() Key {
	return s.keySource.Key()
}

func InMap[K comparable, T any](key K, source BindingSource[T]) BindingSource[T] {
	return &mapBindingSource[K, T]{
		key:       key,
		source:    source,
		keySource: &mapKeySource[K, T]{},
	}
}
