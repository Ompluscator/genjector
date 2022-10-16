package genjector

import (
	"reflect"
	"testing"
)

func Test_baseKeySource_Key(t *testing.T) {
	if new(baseKeySource[int]).Key() != new(baseKeySource[int]).Key() {
		t.Error("keys are different")
	}

	if new(baseKeySource[interface{}]).Key() != new(baseKeySource[interface{}]).Key() {
		t.Error("keys are different")
	}

	if new(baseKeySource[string]).Key() == new(baseKeySource[int]).Key() {
		t.Error("keys are the same")
	}

	if new(baseKeySource[interface{}]).Key() == new(baseKeySource[*int]).Key() {
		t.Error("keys are the same")
	}

	if new(baseKeySource[interface{}]).Key() == new(baseKeySource[*struct{}]).Key() {
		t.Error("keys are the same")
	}

	if new(baseKeySource[struct{}]).Key() == new(baseKeySource[*struct{}]).Key() {
		t.Error("keys are the same")
	}
}

func Test_baseKeySource_Container(t *testing.T) {
	container := map[interface{}]Binding{
		"something": nil,
	}

	first := new(baseKeySource[int]).Container(container)
	second := new(baseKeySource[int]).Container(container)
	if !reflect.DeepEqual(first, second) {
		t.Error("containers are different")
	}

	first = new(baseKeySource[interface{}]).Container(container)
	second = new(baseKeySource[struct{}]).Container(container)
	if !reflect.DeepEqual(first, second) {
		t.Error("containers are different")
	}
}

func Test_sliceKeySource_Key(t *testing.T) {
	if new(sliceKeySource[int]).Key() != new(sliceKeySource[int]).Key() {
		t.Error("keys are different")
	}

	if new(sliceKeySource[interface{}]).Key() != new(sliceKeySource[interface{}]).Key() {
		t.Error("keys are different")
	}

	if new(sliceKeySource[string]).Key() == new(baseKeySource[string]).Key() {
		t.Error("keys are the same")
	}

	if new(sliceKeySource[string]).Key() == new(sliceKeySource[int]).Key() {
		t.Error("keys are the same")
	}

	if new(sliceKeySource[interface{}]).Key() == new(sliceKeySource[*int]).Key() {
		t.Error("keys are the same")
	}

	if new(sliceKeySource[interface{}]).Key() == new(sliceKeySource[*struct{}]).Key() {
		t.Error("keys are the same")
	}

	if new(sliceKeySource[struct{}]).Key() == new(sliceKeySource[*struct{}]).Key() {
		t.Error("keys are the same")
	}
}

func Test_sliceKeySource_Container(t *testing.T) {
	container := map[interface{}]Binding{
		"something": nil,
	}

	first := new(sliceKeySource[int]).Container(container)
	second := new(sliceKeySource[int]).Container(container)
	if !reflect.DeepEqual(first, second) {
		t.Error("containers are different")
	}

	first = new(sliceKeySource[interface{}]).Container(container)
	second = new(sliceKeySource[struct{}]).Container(container)
	if !reflect.DeepEqual(first, second) {
		t.Error("containers are different")
	}
}

func Test_mapKeySource_Key(t *testing.T) {
	if new(mapKeySource[string, int]).Key() != new(mapKeySource[string, int]).Key() {
		t.Error("keys are different")
	}

	if new(mapKeySource[int, interface{}]).Key() != new(mapKeySource[int, interface{}]).Key() {
		t.Error("keys are different")
	}

	if new(mapKeySource[string, string]).Key() == new(baseKeySource[string]).Key() {
		t.Error("keys are the same")
	}

	if new(mapKeySource[string, string]).Key() == new(mapKeySource[string, int]).Key() {
		t.Error("keys are the same")
	}

	if new(mapKeySource[string, string]).Key() == new(mapKeySource[int, string]).Key() {
		t.Error("keys are the same")
	}

	if new(mapKeySource[string, interface{}]).Key() == new(mapKeySource[string, *int]).Key() {
		t.Error("keys are the same")
	}

	if new(mapKeySource[int, interface{}]).Key() == new(mapKeySource[string, *struct{}]).Key() {
		t.Error("keys are the same")
	}

	if new(mapKeySource[string, struct{}]).Key() == new(mapKeySource[string, *struct{}]).Key() {
		t.Error("keys are the same")
	}
}

func Test_mapKeySource_Container(t *testing.T) {
	container := map[interface{}]Binding{
		"something": nil,
	}

	first := new(mapKeySource[string, int]).Container(container)
	second := new(mapKeySource[string, int]).Container(container)
	if !reflect.DeepEqual(first, second) {
		t.Error("containers are different")
	}

	first = new(mapKeySource[int, interface{}]).Container(container)
	second = new(mapKeySource[string, struct{}]).Container(container)
	if !reflect.DeepEqual(first, second) {
		t.Error("containers are different")
	}
}

func Test_sameKeyOption_Key(t *testing.T) {
	key := Key{
		Value:      2,
		Annotation: "value",
	}

	first := new(sameKeyOption).Key(key)
	second := new(sameKeyOption).Key(key)
	if !reflect.DeepEqual(first, second) {
		t.Error("keys are different")
	}
}

func Test_sameKeyOption_Container(t *testing.T) {
	container := map[interface{}]Binding{
		"something": nil,
	}

	first := new(sameKeyOption).Container(container)
	second := new(sameKeyOption).Container(container)
	if !reflect.DeepEqual(first, second) {
		t.Error("containers are different")
	}
}

func Test_annotatedKeyOption_Key(t *testing.T) {
	result := (&annotatedKeyOption{
		annotation: "annotation",
	}).Key(Key{
		Value:      2,
		Annotation: "value",
	})
	if !reflect.DeepEqual(Key{
		Value:      2,
		Annotation: "annotation",
	}, result) {
		t.Error("keys are different")
	}

	result = (&annotatedKeyOption{
		annotation: "",
	}).Key(Key{
		Value:      2,
		Annotation: "value",
	})
	if !reflect.DeepEqual(Key{
		Value:      2,
		Annotation: "",
	}, result) {
		t.Error("keys are different")
	}
}

func Test_annotatedKeyOption_Container(t *testing.T) {
	container := map[interface{}]Binding{
		"something": nil,
	}

	first := new(annotatedKeyOption).Container(container)
	second := new(annotatedKeyOption).Container(container)
	if !reflect.DeepEqual(first, second) {
		t.Error("containers are different")
	}
}

func Test_containerKeyOption_Key(t *testing.T) {
	key := Key{
		Value:      2,
		Annotation: "value",
	}

	first := new(containerKeyOption).Key(key)
	second := new(containerKeyOption).Key(key)
	if !reflect.DeepEqual(first, second) {
		t.Error("keys are different")
	}
}

func Test_containerKeyOption_Container(t *testing.T) {
	result := (&containerKeyOption{
		container: Container{
			"first": nil,
		},
	}).Container(Container{
		"second": nil,
	})
	if !reflect.DeepEqual(Container{
		"first": nil,
	}, result) {
		t.Error("containers are different")
	}

	result = (&containerKeyOption{
		container: nil,
	}).Container(Container{
		"first": nil,
	})
	if result != nil {
		t.Error("containers are different")
	}
}
