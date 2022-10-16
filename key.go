package genjector

type baseKeySource[T any] struct{}

func (baseKeySource[T]) Key() Key {
	return Key{
		Value: (*T)(nil),
	}
}

func (baseKeySource[T]) Container(container Container) Container {
	return container
}

type sliceKeySource[T any] struct{}

func (sliceKeySource[T]) Key() Key {
	return Key{
		Value: (*[]T)(nil),
	}
}

func (sliceKeySource[T]) Container(container Container) Container {
	return container
}

type mapKeySource[K comparable, T any] struct{}

func (mapKeySource[K, T]) Key() Key {
	return Key{
		Value: (*map[K]T)(nil),
	}
}

func (mapKeySource[K, T]) Container(container Container) Container {
	return container
}

type sameKeyOption struct{}

func (sameKeyOption) Key(key Key) Key {
	return key
}

func (sameKeyOption) Container(container Container) Container {
	return container
}

type annotatedKeyOption struct {
	annotation string
}

func (o *annotatedKeyOption) Key(key Key) Key {
	return Key{
		Annotation: o.annotation,
		Value:      key.Value,
	}
}

func (*annotatedKeyOption) Container(container Container) Container {
	return container
}

func AnnotatedWith(annotation string) KeyOption {
	return &annotatedKeyOption{
		annotation: annotation,
	}
}

type containerKeyOption struct {
	container Container
}

func (*containerKeyOption) Key(key Key) Key {
	return key
}

func (o *containerKeyOption) Container(Container) Container {
	return o.container
}

func WithContainer(container Container) KeyOption {
	return &containerKeyOption{
		container: container,
	}
}
