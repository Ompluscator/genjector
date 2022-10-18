package examples

import (
	"github.com/ompluscator/genjector"
	"testing"
)

type SliceInterface interface {
	String() string
}

type SliceStruct struct {
	value string
}

func (s *SliceStruct) Init() {
	s.value = "value provided inside the SliceStruct"
}

func (s *SliceStruct) String() string {
	return s.value
}

func TestAsSlice(t *testing.T) {
	t.Run("Bind multiple pointers to a struct as an implementations for an interface", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(
			genjector.InSlice(genjector.AsPointer[SliceInterface, *SliceStruct]()),
		)
		if err != nil {
			t.Error("binding should not cause an error")
		}

		err = genjector.Bind(
			genjector.InSlice(genjector.AsProvider[SliceInterface](func() (*SliceStruct, error) {
				return &SliceStruct{
					value: "value provided inside the ProviderMethod",
				}, nil
			})),
		)
		if err != nil {
			t.Error("binding should not cause an error")
		}

		err = genjector.Bind(
			genjector.InSlice(genjector.AsInstance[SliceInterface](&SliceStruct{
				value: "value provided inside the Test method",
			})),
		)
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[[]SliceInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance[0].String()
		if value != "value provided inside the SliceStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}

		value = instance[1].String()
		if value != "value provided inside the ProviderMethod" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}

		value = instance[2].String()
		if value != "value provided inside the Test method" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
