package genjector

// baseKeySource is a concrete implementation for KeyOption interface.
type baseKeySource[T any] struct{}

// Key returns the instance of Key that represents a Container key for T type.
//
// It respects KeyOption interface.
func (baseKeySource[T]) Key() Key {
	return Key{
		Value: (*T)(nil),
	}
}

// Container returns the same instance of Container struct provided as an argument.
//
// It respects KeyOption interface.
func (baseKeySource[T]) Container(container Container) Container {
	return container
}

// sliceKeySource is a concrete implementation for KeyOption interface.
type sliceKeySource[T any] struct{}

// Key returns the instance of Key that represents a Container key for slice T types.
//
// It respects KeyOption interface.
func (sliceKeySource[T]) Key() Key {
	return Key{
		Value: (*[]T)(nil),
	}
}

// Container returns the same instance of Container struct provided as an argument.
//
// It respects KeyOption interface.
func (sliceKeySource[T]) Container(container Container) Container {
	return container
}

// mapKeySource is a concrete implementation for KeyOption interface.
type mapKeySource[K comparable, T any] struct{}

// Key returns the instance of Key that represents a Container key for map K-T pairs.
//
// It respects KeyOption interface.
func (mapKeySource[K, T]) Key() Key {
	return Key{
		Value: (*map[K]T)(nil),
	}
}

// Container returns the same instance of Container struct provided as an argument.
//
// It respects KeyOption interface.
func (mapKeySource[K, T]) Container(container Container) Container {
	return container
}

// sameKeyOption is a concrete implementation for KeyOption interface.
type sameKeyOption struct{}

// Key returns the same instance of Key struct provided as an argument.
//
// It respects KeyOption interface.
func (sameKeyOption) Key(key Key) Key {
	return key
}

// Container returns the same instance of Container struct provided as an argument.
//
// It respects KeyOption interface.
func (sameKeyOption) Container(container Container) Container {
	return container
}

// annotatedKeyOption is a concrete implementation for KeyOption interface.
type annotatedKeyOption struct {
	annotation string
}

// Key wrapped instance of Key with a new value for the annotation.
//
// It respects KeyOption interface.
func (o *annotatedKeyOption) Key(key Key) Key {
	return Key{
		Annotation: o.annotation,
		Value:      key.Value,
	}
}

// Container returns the same instance of Container struct provided as an argument.
//
// It respects KeyOption interface.
func (*annotatedKeyOption) Container(container Container) Container {
	return container
}

// containerKeyOption is a concrete implementation for KeyOption interface.
type containerKeyOption struct {
	container Container
}

// Key returns the same instance of Key struct provided as an argument.
//
// It respects KeyOption interface.
func (*containerKeyOption) Key(key Key) Key {
	return key
}

// Container returns an inner instance of Container.
//
// It respects KeyOption interface.
func (o *containerKeyOption) Container(Container) Container {
	return o.container
}
