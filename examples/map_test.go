package examples

import (
	"github.com/ompluscator/genjector"
	"testing"
)

type MapInterface interface {
	String() string
}

type MapStruct struct {
	value string
}

func (s *MapStruct) Init() {
	s.value = "value provided inside the MapStruct"
}

func (s *MapStruct) String() string {
	return s.value
}

func TestAsMap(t *testing.T) {
	t.Run("Bind multiple pointers to a struct as an implementations for an interface", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(
			genjector.InMap("first", genjector.AsPointer[MapInterface, *MapStruct]()),
		)
		if err != nil {
			t.Error("binding should not cause an error")
		}

		err = genjector.Bind(
			genjector.InMap("second", genjector.AsProvider[MapInterface](func() (*MapStruct, error) {
				return &MapStruct{
					value: "value provided inside the ProviderMethod",
				}, nil
			})),
		)
		if err != nil {
			t.Error("binding should not cause an error")
		}

		err = genjector.Bind(
			genjector.InMap("third", genjector.AsInstance[MapInterface](&MapStruct{
				value: "value provided inside the Test method",
			})),
		)
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[map[string]MapInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance["first"].String()
		if value != "value provided inside the MapStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}

		value = instance["second"].String()
		if value != "value provided inside the ProviderMethod" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}

		value = instance["third"].String()
		if value != "value provided inside the Test method" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
