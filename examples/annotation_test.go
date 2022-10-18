package examples

import (
	"fmt"
	"github.com/ompluscator/genjector"
	"testing"
)

type AnnotationInterface interface {
	String() string
}

type AnnotationStruct struct {
	value string
}

func (s *AnnotationStruct) Init() {
	firstChild := genjector.MustNewInstance[*AnnotationChildStruct](genjector.WithAnnotation("first"))
	secondChild := genjector.MustNewInstance[*AnnotationChildStruct](genjector.WithAnnotation("second"))
	s.value = fmt.Sprintf("%s | %s", firstChild.value, secondChild.value)
}

func (s *AnnotationStruct) String() string {
	return s.value
}

type AnnotationChildStruct struct {
	value string
}

func (s *AnnotationChildStruct) Init() {
	s.value = "value provided inside the AnnotationChildStruct"
}

func TestAsAnnotation(t *testing.T) {
	t.Run("Take values from inner 2 child objects defined with proper annotation", func(t *testing.T) {
		genjector.Clean()

		err := genjector.Bind(genjector.AsPointer[AnnotationInterface, *AnnotationStruct]())
		if err != nil {
			t.Error("binding should not cause an error")
		}

		err = genjector.Bind(genjector.AsProvider[*AnnotationChildStruct](func() (*AnnotationChildStruct, error) {
			return &AnnotationChildStruct{
				value: "value from the first child",
			}, nil
		}), genjector.WithAnnotation("first"))
		if err != nil {
			t.Error("binding should not cause an error")
		}

		err = genjector.Bind(genjector.AsProvider[*AnnotationChildStruct](func() (*AnnotationChildStruct, error) {
			return &AnnotationChildStruct{
				value: "value from the second child",
			}, nil
		}), genjector.WithAnnotation("second"))
		if err != nil {
			t.Error("binding should not cause an error")
		}

		instance, err := genjector.NewInstance[AnnotationInterface]()
		if err != nil {
			t.Error("initialization should not cause an error")
		}

		value := instance.String()
		if value != "value from the first child | value from the second child" {
			t.Errorf(`unexpected value received: "%s"`, value)
		}
	})
}
