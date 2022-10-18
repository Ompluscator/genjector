package genjector

import "fmt"

// sliceBinding is a concrete implementation for Binding interface.
type sliceBinding[T any] struct {
	previous *sliceBinding[T]
	current  Binding
}

// Instance returns a slice of T types by executing current Binding and
// all other preceding ones. First it places previous in a slice, and
// then stores the instance of the current.
//
// It respects Binding interface.
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

// sliceBindingSource is a concrete implementation for BindingSource interface.
type sliceBindingSource[T any] struct {
	previous  Binding
	source    BindingSource[T]
	keySource KeySource
}

// Binding returns an instance of a new Binding. If there is no any
// stored predecessor, it will deliver new Binding without containing any
// previous Binding. In case predecessor is defined, all will be returned
// together.
//
// It respects BindingSource interface.
func (b *sliceBindingSource[T]) Binding() (Binding, error) {
	binding, err := b.source.Binding()
	if err != nil {
		return nil, err
	}

	instance, err := binding.Instance(false)
	if err != nil {
		return nil, err
	}

	if _, ok := instance.(T); !ok {
		var initial T
		return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, initial, instance)
	}

	previous, ok := b.previous.(*sliceBinding[T])
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

// SetPrevious stores preceding Binding as a previous one.
//
// It respects FollowingBindingSource interface.
func (b *sliceBindingSource[T]) SetPrevious(binding Binding) {
	b.previous = binding
}

// Key executes the same method from inner KeyOption instance.
//
// It respects BindingOption interface.
func (b *sliceBindingSource[T]) Key() Key {
	return b.keySource.Key()
}

// InSlice delivers a BindingSource for a slice of types T.. It is used as a wrapping
// BindingSource for any other inner. It creates complex Binding in the background
// that stores all T types in a slice and delivers it upon request by executing NewInstance
// method for a slice of T types.
//
// Example:
// err := :genjector.Bind(
//
//	  genjector.InSlice(
//	    genjector.AsInstance[SliceInterface](&SliceStruct{
//	      value: "concrete value",
//		   }),
//	  ),
//
// )
//
// BindingSource can be only used as the first argument to Bind method.
func InSlice[T any](source BindingSource[T]) BindingSource[T] {
	return &sliceBindingSource[T]{
		source:    source,
		keySource: sliceKeySource[T]{},
	}
}

// mapBinding is a concrete implementation for Binding interface.
type mapBinding[K comparable, T any] struct {
	previous *mapBinding[K, T]
	key      K
	current  Binding
}

// Instance returns a map of K-T pairs by executing current Binding and
// all other preceding ones.
//
// It respects Binding interface.
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

// mapBindingSource is a concrete implementation for BindingSource interface.
type mapBindingSource[K comparable, T any] struct {
	previous  Binding
	source    BindingSource[T]
	key       K
	keySource KeySource
}

// Binding returns an instance of a new Binding. If there is no any
// stored predecessor, it will deliver new Binding without containing any
// previous Binding. In case predecessor is defined, all will be returned
// together.
//
// It respects BindingSource interface.
func (b *mapBindingSource[K, T]) Binding() (Binding, error) {
	binding, err := b.source.Binding()
	if err != nil {
		return nil, err
	}

	instance, _ := binding.Instance(false)
	if _, ok := instance.(T); !ok {
		var initial T
		return nil, fmt.Errorf(`binding is not possible for "%v" and "%v"`, initial, instance)
	}

	previous, ok := b.previous.(*mapBinding[K, T])
	if !ok {
		return &mapBinding[K, T]{
			key:     b.key,
			current: binding,
		}, nil
	}

	return &mapBinding[K, T]{
		previous: previous,
		key:      b.key,
		current:  binding,
	}, nil
}

// SetPrevious stores preceding Binding as a previous one.
//
// It respects FollowingBindingSource interface.
func (b *mapBindingSource[K, T]) SetPrevious(binding Binding) {
	b.previous = binding
}

// Key executes the same method from inner KeyOption instance.
//
// It respects BindingOption interface.
func (b *mapBindingSource[K, T]) Key() Key {
	return b.keySource.Key()
}

// InMap delivers a BindingSource for a type T and key's type K, that creates a map
// of K-T pairs. It is used as a wrapping BindingSource for any other inner. It creates
// complex Binding in the background that stores all T types in a map and delivers it
// upon request by executing NewInstance method for a K-T map.
//
// Example:
// err := :genjector.Bind(
//
//	  genjector.InMap(
//	    "third",
//	    genjector.AsInstance[MapInterface](&MapStruct{
//	      value: "concrete value",
//		   }),
//	  ),
//
// )
//
// BindingSource can be only used as the first argument to Bind method.
func InMap[K comparable, T any](key K, source BindingSource[T]) BindingSource[T] {
	return &mapBindingSource[K, T]{
		key:       key,
		source:    source,
		keySource: mapKeySource[K, T]{},
	}
}
