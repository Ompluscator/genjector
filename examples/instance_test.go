package examples

import (
	"testing"

	"github.com/ompluscator/genjector"
)

type InstanceInterface interface {
	String() string
}

type InstanceStruct struct {
	value string
}

func (s *InstanceStruct) Init() {
	s.value = "value provided inside the InstanceStruct"
}

func (s *InstanceStruct) String() string {
	return s.value
}

func TestAsInstance(t *testing.T) {
	t.Run("Bind a pointer to a struct as an implementation for an interface", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsInstance[InstanceInterface](&InstanceStruct{
			value: "value provided inside the Test method",
		}))
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.NewInstance[InstanceInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the Test method" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})

	t.Run("Bind a pointer to a struct as an implementation for the pointer of that struct", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsInstance[*InstanceStruct](&InstanceStruct{
			value: "value provided inside the Test method",
		}))
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.NewInstance[*InstanceStruct]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the Test method" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
