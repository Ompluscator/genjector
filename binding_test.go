package genjector

import (
	"errors"
	"reflect"
	"testing"
)

type testBinding struct {
	value    string
	instance func(initialize bool) (interface{}, error)
}

func (b *testBinding) Instance(initialize bool) (interface{}, error) {
	return b.instance(initialize)
}

type testKeyOption struct {
	key       func(key Key) Key
	container func(container Container) Container
}

func (o *testKeyOption) Key(key Key) Key {
	return o.key(key)
}

func (o *testKeyOption) Container(container Container) Container {
	return o.container(container)
}

type testStruct struct {
	a string
	b int
}

func (s *testStruct) Init() {
	s.a = "test"
	s.b = 10
}

var singleton = 1

type testSingletonStruct struct{}

func (s *testSingletonStruct) Init() {
	singleton++
}

func Test_bindingSource_Binding_error(t *testing.T) {
	binding := &testBinding{}
	binding.instance = func(initialize bool) (interface{}, error) {
		return nil, errors.New("error")
	}

	source := &bindingSource[int]{
		binding: binding,
	}

	result, err := source.Binding()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if result != nil {
		t.Errorf(`expected nil, got: %v`, result)
	}
}

func Test_bindingSource_Binding_wrong(t *testing.T) {
	binding := &testBinding{}
	binding.instance = func(initialize bool) (interface{}, error) {
		return "1", nil
	}

	source := &bindingSource[int]{
		binding: binding,
	}

	result, err := source.Binding()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if result != nil {
		t.Errorf(`expected nil, got: %v`, result)
	}
}

func Test_bindingSource_Binding_success(t *testing.T) {
	binding := &testBinding{}
	binding.instance = func(initialize bool) (interface{}, error) {
		return 1, nil
	}

	source := &bindingSource[int]{
		binding: binding,
	}

	result, err := source.Binding()
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}

	if result == nil {
		t.Error("expected concrete value, got nil")
	}
}

func Test_bindingSource_Key(t *testing.T) {
	source := &bindingSource[int]{
		keySource: baseKeySource[int]{},
	}

	result := source.Key()
	expected := Key{
		Value: (*int)(nil),
	}
	if !reflect.DeepEqual(result, expected) {
		t.Error("expected to get same keys")
	}
}

func Test_valueBinding_Instance_noInitialize(t *testing.T) {
	stringBinding := &valueBinding[string]{}
	result, err := stringBinding.Instance(false)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if result != "" {
		t.Errorf(`expected empty string, got: %v`, result)
	}

	intBinding := &valueBinding[int]{}
	result, err = intBinding.Instance(false)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if result != 0 {
		t.Errorf(`expected 0, got: %v`, result)
	}

	structBinding := &valueBinding[testStruct]{}
	result, err = structBinding.Instance(false)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if !reflect.DeepEqual(result, testStruct{}) {
		t.Errorf(`expected empty struct, got: %v`, result)
	}
}

func Test_valueBinding_Instance_initialize(t *testing.T) {
	stringBinding := &valueBinding[string]{}
	result, err := stringBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if result != "" {
		t.Errorf(`expected empty string, got: %v`, result)
	}

	intBinding := &valueBinding[int]{}
	result, err = intBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if result != 0 {
		t.Errorf(`expected 0, got: %v`, result)
	}

	structBinding := &valueBinding[testStruct]{}
	result, err = structBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if !reflect.DeepEqual(result, testStruct{
		a: "test",
		b: 10,
	}) {
		t.Errorf(`expected empty struct, got: %v`, result)
	}
}

func TestAsValue(t *testing.T) {
	result := AsValue[string, string]()
	if !reflect.DeepEqual(result, &bindingSource[string]{
		binding:   valueBinding[string]{},
		keySource: baseKeySource[string]{},
	}) {
		t.Error("binding sources are different")
	}

	result = AsValue[interface{}, testStruct]()
	if !reflect.DeepEqual(result, &bindingSource[interface{}]{
		binding:   valueBinding[testStruct]{},
		keySource: baseKeySource[interface{}]{},
	}) {
		t.Error("binding sources are different")
	}
}

func Test_referenceBinding_Instance_noInitialize(t *testing.T) {
	stringBinding := &referenceBinding[string]{}
	result, err := stringBinding.Instance(false)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	stringActual := result.(*string)
	if *stringActual != "" {
		t.Errorf(`expected reference to empty string, got: %v`, result)
	}

	intBinding := &referenceBinding[int]{}
	result, err = intBinding.Instance(false)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	intActual := result.(*int)
	if *intActual != 0 {
		t.Errorf(`expected reference to 0, got: %v`, result)
	}

	structBinding := &referenceBinding[testStruct]{}
	result, err = structBinding.Instance(false)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if !reflect.DeepEqual(result, &testStruct{}) {
		t.Errorf(`expected reference to empty struct, got: %v`, result)
	}
}

func Test_referenceBinding_Instance_initialize(t *testing.T) {
	stringBinding := &referenceBinding[string]{}
	result, err := stringBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	stringActual := result.(*string)
	if *stringActual != "" {
		t.Errorf(`expected reference to empty string, got: %v`, result)
	}

	intBinding := &referenceBinding[int]{}
	result, err = intBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	intActual := result.(*int)
	if *intActual != 0 {
		t.Errorf(`expected reference to 0, got: %v`, result)
	}

	structBinding := &referenceBinding[testStruct]{}
	result, err = structBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if !reflect.DeepEqual(result, &testStruct{
		a: "test",
		b: 10,
	}) {
		t.Errorf(`expected reference to empty struct, got: %v`, result)
	}
}

func TestAsReference(t *testing.T) {
	result := AsReference[*string, *string]()
	if !reflect.DeepEqual(result, &bindingSource[*string]{
		binding:   referenceBinding[string]{},
		keySource: baseKeySource[*string]{},
	}) {
		t.Error("binding sources are different")
	}

	result = AsReference[interface{}, *testStruct]()
	if !reflect.DeepEqual(result, &bindingSource[interface{}]{
		binding:   referenceBinding[testStruct]{},
		keySource: baseKeySource[interface{}]{},
	}) {
		t.Error("binding sources are different")
	}
}

func Test_ProviderMethod_Instance(t *testing.T) {
	var stringBinding ProviderMethod[string] = func() (string, error) {
		return "value", nil
	}
	result, err := stringBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if result != "value" {
		t.Errorf(`expected reference to empty string, got: %v`, result)
	}

	var intBinding ProviderMethod[*int] = func() (*int, error) {
		value := 10
		return &value, nil
	}
	result, err = intBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	intActual := result.(*int)
	if *intActual != 10 {
		t.Errorf(`expected reference to 0, got: %v`, result)
	}

	var structBinding ProviderMethod[*testStruct] = func() (*testStruct, error) {
		return &testStruct{
			a: "a",
			b: 5,
		}, nil
	}
	result, err = structBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if !reflect.DeepEqual(result, &testStruct{
		a: "a",
		b: 5,
	}) {
		t.Errorf(`expected reference to empty struct, got: %v`, result)
	}
}

func TestAsProvider(t *testing.T) {
	result := AsProvider[*testStruct](func() (*testStruct, error) {
		return &testStruct{
			a: "a",
			b: 5,
		}, nil
	})
	binding, err := result.Binding()
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}

	if !reflect.DeepEqual(instance, &testStruct{
		a: "a",
		b: 5,
	}) {
		t.Error("instances are different")
	}
}

func Test_instanceBinding_Instance(t *testing.T) {
	stringBinding := &instanceBinding[string]{
		instance: "value",
	}
	result, err := stringBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if result != "value" {
		t.Errorf(`expected reference to empty string, got: %v`, result)
	}

	intBinding := &instanceBinding[int]{
		instance: 10,
	}
	result, err = intBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if result != 10 {
		t.Errorf(`expected reference to 0, got: %v`, result)
	}

	structBinding := &instanceBinding[*testStruct]{
		instance: &testStruct{
			a: "a",
			b: 5,
		},
	}
	result, err = structBinding.Instance(true)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if !reflect.DeepEqual(result, &testStruct{
		a: "a",
		b: 5,
	}) {
		t.Errorf(`expected reference to empty struct, got: %v`, result)
	}
}

func TestAsInstance(t *testing.T) {
	result := AsInstance[string]("value")
	if !reflect.DeepEqual(result, &bindingSource[string]{
		binding: &instanceBinding[string]{
			instance: "value",
		},
		keySource: baseKeySource[string]{},
	}) {
		t.Error("binding sources are different")
	}

	result = AsInstance[interface{}](&testStruct{
		a: "test",
	})
	if !reflect.DeepEqual(result, &bindingSource[interface{}]{
		binding: &instanceBinding[*testStruct]{
			instance: &testStruct{
				a: "test",
			},
		},
		keySource: baseKeySource[interface{}]{},
	}) {
		t.Error("binding sources are different")
	}
}

func Test_bindingOption_Binding(t *testing.T) {
	option := &bindingOption{
		bindingFunc: func(binding Binding) (Binding, error) {
			return &instanceBinding[string]{
				instance: "value",
			}, nil
		},
	}

	binding, err := option.Binding(nil)
	if err != nil {
		t.Errorf(`expected nil, got: %v`, err)
	}
	if !reflect.DeepEqual(binding, &instanceBinding[string]{
		instance: "value",
	}) {
		t.Error("bindings are different")
	}
}

func Test_bindingOption_Container(t *testing.T) {
	option := &bindingOption{
		keyOption: &testKeyOption{
			container: func(container Container) Container {
				return Container{
					"something": nil,
				}
			},
		},
	}

	container := option.Container(nil)
	if !reflect.DeepEqual(container, Container{
		"something": nil,
	}) {
		t.Error("containers are different")
	}
}

func Test_bindingOption_Key(t *testing.T) {
	option := &bindingOption{
		keyOption: &testKeyOption{
			key: func(key Key) Key {
				return Key{
					Annotation: "annotation",
				}
			},
		},
	}

	key := option.Key(Key{})
	if !reflect.DeepEqual(key, Key{
		Annotation: "annotation",
	}) {
		t.Error("containers are different")
	}
}

func Test_singletonBinding_Instance(t *testing.T) {
	option := &bindingOption{
		bindingFunc: func(binding Binding) (Binding, error) {
			return &singletonBinding{
				parent: binding,
			}, nil
		},
	}

	binding, err := option.Binding(&referenceBinding[testSingletonStruct]{})
	if err != nil {
		t.Error("unexpected error")
	}

	instance, err := binding.Instance(true)
	if err != nil {
		t.Error("unexpected error")
	}
	if !reflect.DeepEqual(instance, &testSingletonStruct{}) {
		t.Error("instance are different")
	}
	if singleton != 2 {
		t.Errorf(`expected value 2, go: %d`, singleton)
	}

	instance, err = binding.Instance(true)
	if err != nil {
		t.Error("unexpected error")
	}
	if !reflect.DeepEqual(instance, &testSingletonStruct{}) {
		t.Error("instance are different")
	}
	if singleton != 2 {
		t.Errorf(`expected value 2, go: %d`, singleton)
	}
}

func TestAsSingleton(t *testing.T) {
	result := AsSingleton()

	binding, err := result.(*bindingOption).bindingFunc(&valueBinding[int]{})
	if err != nil {
		t.Error("unexpected error")
	}
	if !reflect.DeepEqual(binding, &singletonBinding{
		parent: &valueBinding[int]{},
	}) {
		t.Error("bindings are different")
	}

	result.(*bindingOption).bindingFunc = nil

	if !reflect.DeepEqual(result, &bindingOption{
		keyOption: sameKeyOption{},
	}) {
		t.Error("binding options are different")
	}
}

func TestWithAnnotation(t *testing.T) {
	result := WithAnnotation("annotation")
	result.(*bindingOption).bindingFunc = nil
	if !reflect.DeepEqual(result, &bindingOption{
		keyOption: &annotatedKeyOption{
			annotation: "annotation",
		},
	}) {
		t.Error("binding options are different")
	}
}

func TestWithContainer(t *testing.T) {
	result := WithContainer(Container{
		"first": nil,
	})
	if !reflect.DeepEqual(&containerKeyOption{
		container: Container{
			"first": nil,
		},
	}, result.(*bindingOption).keyOption) {
		t.Error("containerKeyOption does not contain the right value")
	}

	result = WithContainer(nil)
	if !reflect.DeepEqual(new(containerKeyOption), result.(*bindingOption).keyOption) {
		t.Error("containerKeyOption does not contain the right value")
	}
}
