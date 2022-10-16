package examples

import (
	"github.com/ompluscator/genjector"
	"testing"
)

type ReferenceInterface interface {
	String() string
}

type ReferenceStruct struct {
	value string
}

func (s *ReferenceStruct) Init() {
	s.value = "value provided inside the ReferenceStruct"
}

func (s *ReferenceStruct) String() string {
	return s.value
}

func TestAsReference(t *testing.T) {
	t.Run("Bind a reference to a struct as an implementation for an interface", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsReference[ReferenceInterface, *ReferenceStruct]())
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[ReferenceInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the ReferenceStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})

	t.Run("Bind a reference to a struct as an implementation for the reference of that struct", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsReference[*ReferenceStruct, *ReferenceStruct]())
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[*ReferenceStruct]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the ReferenceStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
