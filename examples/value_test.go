package examples

import (
	"github.com/ompluscator/genjector"
	"testing"
)

type ValueInterface interface {
	String() string
}

type ValueStruct struct {
	value string
}

func (s *ValueStruct) Init() {
	s.value = "value provided inside the ValueStruct"
}

func (s ValueStruct) String() string {
	return s.value
}

func TestAsValue(t *testing.T) {
	t.Run("Bind a struct as an implementation for an interface", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsValue[ValueInterface, ValueStruct]())
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[ValueInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the ValueStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})

	t.Run("Bind a struct as an implementation for that struct", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsValue[ValueStruct, ValueStruct]())
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[ValueStruct]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the ValueStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
