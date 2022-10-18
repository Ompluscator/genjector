package genjector

import (
	"errors"
	"reflect"
	"testing"
)

func TestKey_Generate(t *testing.T) {
	key := Key{}
	generated := key.Generate()
	if generated != nil {
		t.Errorf("expected nil, got %v", generated)
	}

	key = Key{
		Annotation: "annotation",
	}
	generated = key.Generate()
	if !reflect.DeepEqual(generated, [2]interface{}{"annotation", nil}) {
		t.Errorf("expected concrete value, got %v", generated)
	}

	key = Key{
		Annotation: "annotation",
		Value:      "value",
	}
	generated = key.Generate()
	if !reflect.DeepEqual(generated, [2]interface{}{"annotation", "value"}) {
		t.Errorf("expected concrete value, got %v", generated)
	}
}

func TestNewContainer(t *testing.T) {
	container := NewContainer()
	if !reflect.DeepEqual(container, Container{}) {
		t.Errorf("expected concrete value, got %v", container)
	}
}

func TestBind_bindingError(t *testing.T) {
	err := Bind[int](&testBindingSource{
		binding: func() (Binding, error) {
			return nil, errors.New("error")
		},
		key: func() Key {
			return Key{}
		},
	})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestBind_bindingOptionError(t *testing.T) {
	err := Bind[int](&testBindingSource{
		binding: func() (Binding, error) {
			return &testBinding{}, nil
		},
		key: func() Key {
			return Key{}
		},
	}, &testBindingOption{
		binding: func(binding Binding) (Binding, error) {
			return nil, errors.New("error")
		},
		key: func(key Key) Key {
			return key
		},
		container: func(container Container) Container {
			return NewContainer()
		},
	})
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestBind_binding_success(t *testing.T) {
	inner := NewContainer()
	if !reflect.DeepEqual(inner, Container{}) {
		t.Errorf("expected concrete value, got %v", inner)
	}

	err := Bind[int](&testBindingSource{
		binding: func() (Binding, error) {
			return &testBinding{}, nil
		},
		key: func() Key {
			return Key{
				Annotation: "first",
				Value:      "1",
			}
		},
	}, &testBindingOption{
		binding: func(binding Binding) (Binding, error) {
			binding.(*testBinding).value = "value"
			return binding, nil
		},
		key: func(key Key) Key {
			return Key{
				Annotation: key.Annotation + "second",
				Value:      key.Value.(string) + "2",
			}
		},
		container: func(Container) Container {
			return inner
		},
	})
	if err != nil {
		t.Errorf("expected nil, got error %v", err)
	}

	if !reflect.DeepEqual(inner, Container{
		[2]interface{}{
			"firstsecond",
			"12",
		}: &testBinding{
			value: "value",
		},
	}) {
		t.Errorf("expected concrete value, got %v", inner)
	}
}

func TestMustBind(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("the code did not panic")
		}
	}()

	MustBind[int](&testBindingSource{
		binding: func() (Binding, error) {
			return nil, errors.New("error")
		},
		key: func() Key {
			return Key{}
		},
	})
}

func TestNewInstance_defaultValue(t *testing.T) {
	instance, err := NewInstance[testStruct]()
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, testStruct{
		a: "test",
		b: 10,
	}) {
		t.Errorf("expected concrete value, got %v", instance)
	}
}

func TestNewInstance_defaultPointer(t *testing.T) {
	instance, err := NewInstance[*testStruct]()
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if instance != nil {
		t.Errorf("expected nil, got %v", instance)
	}
}

func TestNewInstance_defaultError(t *testing.T) {
	instance, err := NewInstance[interface{}]()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf("expected nil, got %v", instance)
	}
}

func TestNewInstance_success(t *testing.T) {
	inner := Container{
		(*int)(nil): &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return nil, errors.New("error")
			},
		},
	}

	instance, err := NewInstance[int](WithContainer(inner))
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != 0 {
		t.Errorf("expected 0, got %v", instance)
	}
}

func TestNewInstance_invalid(t *testing.T) {
	inner := Container{
		(*int)(nil): &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return "value", nil
			},
		},
	}

	instance, err := NewInstance[int](WithContainer(inner))
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != 0 {
		t.Errorf("expected 0, got %v", instance)
	}
}

func TestNewInstance_simple_success(t *testing.T) {
	inner := Container{
		(*int)(nil): &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return 10, nil
			},
		},
	}

	instance, err := NewInstance[int](WithContainer(inner))
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if instance != 10 {
		t.Errorf("expected 10, got %v", instance)
	}
}

func TestNewInstance_complex_success(t *testing.T) {
	inner := Container{
		[2]interface{}{"annotation", nil}: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return 10, nil
			},
		},
	}

	instance, err := NewInstance[int](&testKeyOption{
		key: func(key Key) Key {
			return Key{
				Annotation: "annotation",
			}
		},
		container: func(container Container) Container {
			return inner
		},
	})
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if instance != 10 {
		t.Errorf("expected 10, got %v", instance)
	}
}
