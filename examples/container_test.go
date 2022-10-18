package examples

import (
	"github.com/ompluscator/genjector"
	"testing"
)

type ContainerInterface interface {
	String() string
}

type ContainerStruct struct {
	value string
}

func (s *ContainerStruct) Init() {
	s.value = "value provided inside the ContainerStruct"
}

func (s *ContainerStruct) String() string {
	return s.value
}

func TestAsContainer(t *testing.T) {
	t.Run("Store binding inside the customer container and not the global one", func(t *testing.T) {
		genjector.Clean()

		customContainer := genjector.NewContainer()

		err := genjector.Bind(
			genjector.AsPointer[ContainerInterface, *ContainerStruct](),
			genjector.WithContainer(customContainer),
		)
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.Initialize[ContainerInterface]()
		if err == nil {
			t.Error("expected an error, but got nil")
		}
		if instance != nil {
			t.Errorf(`unexpected instance received: "%s"`, instance)
		}

		instance, err = genjector.Initialize[ContainerInterface](genjector.WithContainer(customContainer))
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the ContainerStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
