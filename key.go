package genjector

type baseKeySource[T any] struct{}

func (*baseKeySource[T]) Key() Key {
	return Key{
		Value: (*T)(nil),
	}
}

func (*baseKeySource[T]) Container(container map[interface{}]Binding) map[interface{}]Binding {
	return container
}

type sliceKeySource[T any] struct{}

func (*sliceKeySource[T]) Key() Key {
	return Key{
		Value: (*[]T)(nil),
	}
}

func (*sliceKeySource[T]) Container(container map[interface{}]Binding) map[interface{}]Binding {
	return container
}

type mapKeySource[K comparable, T any] struct{}

func (*mapKeySource[K, T]) Key() Key {
	return Key{
		Value: (*map[K]T)(nil),
	}
}

func (*mapKeySource[K, T]) Container(container map[interface{}]Binding) map[interface{}]Binding {
	return container
}

type sameKeyOption struct{}

func (*sameKeyOption) Key(key Key) Key {
	return key
}

func (*sameKeyOption) Container(container map[interface{}]Binding) map[interface{}]Binding {
	return container
}

type annotatedKeyOption struct {
	annotation string
}

func (o *annotatedKeyOption) Key(Key) Key {
	return Key{
		Annotation: o.annotation,
	}
}

func (*annotatedKeyOption) Container(container map[interface{}]Binding) map[interface{}]Binding {
	return container
}

func AnnotatedWith(annotation string) KeyOption {
	return &annotatedKeyOption{
		annotation: annotation,
	}
}

type containerKeyOption struct {
	container map[interface{}]Binding
}

func (*containerKeyOption) Key(key Key) Key {
	return key
}

func (o *containerKeyOption) Container(map[interface{}]Binding) map[interface{}]Binding {
	return o.container
}

func WithContainer(container map[interface{}]Binding) KeyOption {
	return &containerKeyOption{
		container: container,
	}
}
