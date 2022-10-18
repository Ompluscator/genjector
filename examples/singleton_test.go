package examples

import (
	"github.com/ompluscator/genjector"
	"testing"
)

var counter = 0

type SingletonInterface interface {
	String() string
}

type SingletonStruct struct {
	value string
}

func (s *SingletonStruct) Init() {
	counter++
	s.value = "value provided inside the SingletonStruct"
}

func (s *SingletonStruct) String() string {
	return s.value
}

func TestAsSingleton(t *testing.T) {
	t.Run("Expecting counter to stay at 1 when instance is defined as a singleton", func(t *testing.T) {
		genjector.Clean()
		counter = 0

		err := genjector.Bind(
			genjector.AsPointer[SingletonInterface, *SingletonStruct](),
			genjector.AsSingleton(),
		)
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.NewInstance[SingletonInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the SingletonStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
		if counter != 1 {
			t.Errorf(`unexpected counter received: "%d"`, counter)
		}

		instance, err = genjector.NewInstance[SingletonInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value = instance.String()
		if value != "value provided inside the SingletonStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
		if counter != 1 {
			t.Errorf(`unexpected counter received: "%d"`, counter)
		}
	})

	t.Run("Expecting counter to be greater than 1 when instance is not defined as a singleton", func(t *testing.T) {
		genjector.Clean()
		counter = 0

		err := genjector.Bind(genjector.AsPointer[SingletonInterface, *SingletonStruct]())
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.NewInstance[SingletonInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value provided inside the SingletonStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
		if counter != 1 {
			t.Errorf(`unexpected counter received: "%d"`, counter)
		}

		instance, err = genjector.NewInstance[SingletonInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value = instance.String()
		if value != "value provided inside the SingletonStruct" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
		if counter <= 1 {
			t.Errorf(`unexpected counter received: "%d"`, counter)
		}
	})
}
