package genjector

import (
	"errors"
	"reflect"
	"testing"
)

type testBindingSource struct {
	binding func() (Binding, error)
	key     func() Key
}

func (o *testBindingSource) Binding() (Binding, error) {
	return o.binding()
}

func (o *testBindingSource) Key() Key {
	return o.key()
}

type testBindingOption struct {
	binding   func(binding Binding) (Binding, error)
	key       func(key Key) Key
	container func(container Container) Container
}

func (o *testBindingOption) Binding(binding Binding) (Binding, error) {
	return o.binding(binding)
}

func (o *testBindingOption) Key(key Key) Key {
	return o.key(key)
}

func (o *testBindingOption) Container(container Container) Container {
	return o.container(container)
}

func Test_sliceBinding_firstItem_error(t *testing.T) {
	binding := &sliceBinding[testStruct]{
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return nil, errors.New("error")
			},
		},
	}

	instance, err := binding.Instance(true)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf(`expected nil, go %v`, instance)
	}
}

func Test_sliceBinding_firstItem_invalid(t *testing.T) {
	binding := &sliceBinding[testStruct]{
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return "value", nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf(`expected nil, go %v`, instance)
	}
}

func Test_sliceBinding_firstItem_success(t *testing.T) {
	binding := &sliceBinding[testStruct]{
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, []testStruct{
		{
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_sliceBinding_slice_error(t *testing.T) {
	binding := &sliceBinding[testStruct]{
		previous: &sliceBinding[testStruct]{
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return nil, errors.New("error")
				},
			},
		},
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf(`expected nil, go %v`, instance)
	}
}

func Test_sliceBinding_slice_invalid(t *testing.T) {
	binding := &sliceBinding[testStruct]{
		previous: &sliceBinding[testStruct]{
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return "value", nil
				},
			},
		},
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf(`expected nil, go %v`, instance)
	}
}

func Test_sliceBinding_slice_success(t *testing.T) {
	binding := &sliceBinding[testStruct]{
		previous: &sliceBinding[testStruct]{
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return testStruct{
						a: "first",
						b: 5,
					}, nil
				},
			},
		},
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, []testStruct{
		{
			a: "first",
			b: 5,
		},
		{
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_sliceBinding_slice_success_noInitialize(t *testing.T) {
	binding := &sliceBinding[testStruct]{
		previous: &sliceBinding[testStruct]{
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return testStruct{
						a: "first",
						b: 5,
					}, nil
				},
			},
		},
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(false)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, []testStruct{
		{
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_sliceBindingSource_Binding_error(t *testing.T) {
	source := &sliceBindingSource[testStruct]{
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return nil, errors.New("error")
			},
		},
	}

	binding, err := source.Binding()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if binding != nil {
		t.Errorf(`expected nil, go %v`, binding)
	}
}

func Test_sliceBindingSource_Binding_instanceError(t *testing.T) {
	source := &sliceBindingSource[testStruct]{
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return &testBinding{
					instance: func(initialize bool) (interface{}, error) {
						return nil, errors.New("error")
					},
				}, nil
			},
		},
	}

	binding, err := source.Binding()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if binding != nil {
		t.Errorf(`expected nil, go %v`, binding)
	}
}

func Test_sliceBindingSource_Binding_invalid(t *testing.T) {
	source := &sliceBindingSource[testStruct]{
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return &testBinding{
					instance: func(initialize bool) (interface{}, error) {
						return "value", nil
					},
				}, nil
			},
		},
	}

	binding, err := source.Binding()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if binding != nil {
		t.Errorf(`expected nil, go %v`, binding)
	}
}

func Test_sliceBindingSource_Binding_current_success(t *testing.T) {
	source := &sliceBindingSource[testStruct]{
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return &testBinding{
					instance: func(initialize bool) (interface{}, error) {
						return testStruct{
							a: "value",
							b: 20,
						}, nil
					},
				}, nil
			},
		},
	}

	binding, err := source.Binding()
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, []testStruct{
		{
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_sliceBindingSource_Binding_previous_success(t *testing.T) {
	source := &sliceBindingSource[testStruct]{
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return &testBinding{
					instance: func(initialize bool) (interface{}, error) {
						return testStruct{
							a: "value",
							b: 20,
						}, nil
					},
				}, nil
			},
		},
		previous: &sliceBinding[testStruct]{
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return testStruct{
						a: "first",
						b: 5,
					}, nil
				},
			},
		},
	}

	binding, err := source.Binding()
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, []testStruct{
		{
			a: "first",
			b: 5,
		},
		{
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_sliceBindingSource_SetPrevious(t *testing.T) {
	source := &sliceBindingSource[int]{}
	source.SetPrevious(&testBinding{})

	if !reflect.DeepEqual(source, &sliceBindingSource[int]{
		previous: &testBinding{},
	}) {
		t.Error("expected source to match concrete value")
	}
}

func Test_sliceBindingSource_Key(t *testing.T) {
	source := &sliceBindingSource[int]{
		keySource: sliceKeySource[int]{},
	}
	key := source.Key()

	if !reflect.DeepEqual(key, Key{
		Value: (*[]int)(nil),
	}) {
		t.Error("expected source to match concrete value")
	}
}

func TestInSlice(t *testing.T) {
	source := InSlice[string](&testBindingSource{})

	if !reflect.DeepEqual(source, &sliceBindingSource[string]{
		source:    &testBindingSource{},
		keySource: sliceKeySource[string]{},
	}) {
		t.Error("expected source to match concrete value")
	}
}

func Test_mapBinding_firstItem_error(t *testing.T) {
	binding := &mapBinding[string, testStruct]{
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return nil, errors.New("error")
			},
		},
	}

	instance, err := binding.Instance(true)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf(`expected nil, go %v`, instance)
	}
}

func Test_mapBinding_firstItem_invalid(t *testing.T) {
	binding := &mapBinding[string, testStruct]{
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return "value", nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf(`expected nil, go %v`, instance)
	}
}

func Test_mapBinding_firstItem_success(t *testing.T) {
	binding := &mapBinding[string, testStruct]{
		key: "value",
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, map[string]testStruct{
		"value": {
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_mapBinding_map_error(t *testing.T) {
	binding := &mapBinding[string, testStruct]{
		previous: &mapBinding[string, testStruct]{
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return nil, errors.New("error")
				},
			},
		},
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf(`expected nil, go %v`, instance)
	}
}

func Test_mapBinding_map_invalid(t *testing.T) {
	binding := &mapBinding[string, testStruct]{
		previous: &mapBinding[string, testStruct]{
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return "value", nil
				},
			},
		},
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if instance != nil {
		t.Errorf(`expected nil, go %v`, instance)
	}
}

func Test_mapBinding_map_success(t *testing.T) {
	binding := &mapBinding[string, testStruct]{
		previous: &mapBinding[string, testStruct]{
			key: "first",
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return testStruct{
						a: "first",
						b: 5,
					}, nil
				},
			},
		},
		key: "value",
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, map[string]testStruct{
		"first": {
			a: "first",
			b: 5,
		},
		"value": {
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_mapBinding_map_success_noInitialize(t *testing.T) {
	binding := &mapBinding[string, testStruct]{
		previous: &mapBinding[string, testStruct]{
			key: "first",
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return testStruct{
						a: "first",
						b: 5,
					}, nil
				},
			},
		},
		key: "value",
		current: &testBinding{
			instance: func(initialize bool) (interface{}, error) {
				return testStruct{
					a: "value",
					b: 20,
				}, nil
			},
		},
	}

	instance, err := binding.Instance(false)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, map[string]testStruct{
		"value": {
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_mapBindingSource_Binding_error(t *testing.T) {
	source := &mapBindingSource[string, testStruct]{
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return nil, errors.New("error")
			},
		},
	}

	binding, err := source.Binding()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if binding != nil {
		t.Errorf(`expected nil, go %v`, binding)
	}
}

func Test_mapBindingSource_Binding_instanceError(t *testing.T) {
	source := &mapBindingSource[string, testStruct]{
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return &testBinding{
					instance: func(initialize bool) (interface{}, error) {
						return nil, errors.New("error")
					},
				}, nil
			},
		},
	}

	binding, err := source.Binding()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if binding != nil {
		t.Errorf(`expected nil, go %v`, binding)
	}
}

func Test_mapBindingSource_Binding_invalid(t *testing.T) {
	source := &mapBindingSource[string, testStruct]{
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return &testBinding{
					instance: func(initialize bool) (interface{}, error) {
						return "value", nil
					},
				}, nil
			},
		},
	}

	binding, err := source.Binding()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if binding != nil {
		t.Errorf(`expected nil, go %v`, binding)
	}
}

func Test_mapBindingSource_Binding_current_success(t *testing.T) {
	source := &mapBindingSource[string, testStruct]{
		key: "value",
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return &testBinding{
					instance: func(initialize bool) (interface{}, error) {
						return testStruct{
							a: "value",
							b: 20,
						}, nil
					},
				}, nil
			},
		},
	}

	binding, err := source.Binding()
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, map[string]testStruct{
		"value": {
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_mapBindingSource_Binding_previous_success(t *testing.T) {
	source := &mapBindingSource[string, testStruct]{
		key: "value",
		source: &testBindingSource{
			binding: func() (Binding, error) {
				return &testBinding{
					instance: func(initialize bool) (interface{}, error) {
						return testStruct{
							a: "value",
							b: 20,
						}, nil
					},
				}, nil
			},
		},
		previous: &mapBinding[string, testStruct]{
			key: "first",
			current: &testBinding{
				instance: func(initialize bool) (interface{}, error) {
					return testStruct{
						a: "first",
						b: 5,
					}, nil
				},
			},
		},
	}

	binding, err := source.Binding()
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf("expected nil, got error %s", err)
	}

	if !reflect.DeepEqual(instance, map[string]testStruct{
		"first": {
			a: "first",
			b: 5,
		},
		"value": {
			a: "value",
			b: 20,
		},
	}) {
		t.Error("expected instance to match concrete value")
	}
}

func Test_mapBindingSource_SetPrevious(t *testing.T) {
	source := &mapBindingSource[string, int]{}
	source.SetPrevious(&testBinding{})

	if !reflect.DeepEqual(source, &mapBindingSource[string, int]{
		previous: &testBinding{},
	}) {
		t.Error("expected source to match concrete value")
	}
}

func Test_mapBindingSource_Key(t *testing.T) {
	source := &mapBindingSource[string, int]{
		keySource: mapKeySource[string, int]{},
	}
	key := source.Key()

	if !reflect.DeepEqual(key, Key{
		Value: (*map[string]int)(nil),
	}) {
		t.Error("expected source to match concrete value")
	}
}

func TestInMap(t *testing.T) {
	source := InMap[string, int]("key", &testBindingSource{})

	if !reflect.DeepEqual(source, &mapBindingSource[string, int]{
		source:    &testBindingSource{},
		key:       "key",
		keySource: mapKeySource[string, int]{},
	}) {
		t.Error("expected source to match concrete value")
	}
}
