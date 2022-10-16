package examples

import (
	"testing"

	"github.com/ompluscator/genjector"
)

type ProviderInterface interface {
	String() string
}

type ProviderStruct struct {
	value string
}

func (s *ProviderStruct) Init() {
	s.value = "value provided inside the ProviderStruct"
}

func (s *ProviderStruct) String() string {
	return s.value
}

func TestAsProvider(t *testing.T) {
	t.Run("Bind a ProviderMethod for a struct as an implementation for an interface", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsProvider[ProviderInterface](func() (*ProviderStruct, error) {
			return &ProviderStruct{
				value: "value provided inside the ProviderMethod",
			}, nil
		}))
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[ProviderInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the ProviderMethod" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})

	t.Run("Bind a ProviderMethod for a struct as an implementation for that struct", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsProvider[*ProviderStruct](func() (*ProviderStruct, error) {
			return &ProviderStruct{
				value: "value provided inside the ProviderMethod",
			}, nil
		}))

		instance, err := genjector.Initialize[*ProviderStruct]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the ProviderMethod" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
