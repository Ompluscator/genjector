package examples

import (
	"github.com/ompluscator/genjector"
	"testing"
)

type PointerInterface interface {
	String() string
}

type PointerStruct struct {
	value string
}

func (s *PointerStruct) Init() {
	s.value = "value provided inside the PointerStruct"
}

func (s *PointerStruct) String() string {
	return s.value
}

func TestAsPointer(t *testing.T) {
	t.Run("Bind a pointer to a struct as an implementation for an interface", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsPointer[PointerInterface, *PointerStruct]())
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[PointerInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the PointerStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})

	t.Run("Bind a pointer to a struct as an implementation for the pointer of that struct", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsPointer[*PointerStruct, *PointerStruct]())
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[*PointerStruct]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the PointerStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
