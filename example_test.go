package genjector_test

import (
	"github.com/ompluscator/genjector"
	"strconv"
	"testing"
)

type ExampleInterface interface {
	First() string
	Second(arg string) int
}

type ExampleOne struct {
	first  string
	second int
}

func (e *ExampleOne) Init() {
	e.first = "example"
	e.second = 10000
}

func (e *ExampleOne) First() string {
	return e.first
}

func (e *ExampleOne) Second(arg string) int {
	value, _ := strconv.ParseInt(arg, 10, 64)
	return e.second + int(value)
}

func TestExampleOne(t *testing.T) {
	genjector.Clean()

	err := genjector.Bind(genjector.AsReference[ExampleInterface, *ExampleOne]())
	if err != nil {
		t.Error("binding should not cause an error")
	}

	instance, err := genjector.Initialize[ExampleInterface]()
	if err != nil {
		t.Error("initialization should not cause an error")
	}

	if instance.First() != "example" {
		t.Error("unexpected value from method First()")
	}

	if instance.Second("1") != 10001 {
		t.Error(`unexpected value from method Second("1")`)
	}

	err = genjector.Bind(genjector.AsReference[*ExampleOne, *ExampleOne]())
	if err != nil {
		t.Error("binding should not cause an error")
	}

	instance, err = genjector.Initialize[*ExampleOne]()
	if err != nil {
		t.Error("initialization should not cause an error")
	}

	if instance.First() != "example" {
		t.Error("unexpected value from method First()")
	}

	if instance.Second("1") != 10001 {
		t.Error(`unexpected value from method Second("1")`)
	}
}

type ExampleTwo struct{}

func (e ExampleTwo) First() string {
	return "example2"
}

func (e ExampleTwo) Second(arg string) int {
	value, _ := strconv.ParseInt(arg, 10, 64)
	return 20000 + int(value)
}

func TestExampleTwo(t *testing.T) {
	genjector.Clean()

	err := genjector.Bind(genjector.AsValue[ExampleInterface, ExampleTwo]())
	if err != nil {
		t.Error("binding should not cause an error")
	}

	instance, err := genjector.Initialize[ExampleInterface]()
	if err != nil {
		t.Error("initialization should not cause an error")
	}

	if instance.First() != "example2" {
		t.Error("unexpected value from method First()")
	}

	if instance.Second("1") != 20001 {
		t.Error(`unexpected value from method Second("1")`)
	}

	err = genjector.Bind(genjector.AsValue[ExampleTwo, ExampleTwo]())
	if err != nil {
		t.Error("binding should not cause an error")
	}

	instance, err = genjector.Initialize[ExampleTwo]()
	if err != nil {
		t.Error("initialization should not cause an error")
	}

	if instance.First() != "example2" {
		t.Error("unexpected value from method First()")
	}

	if instance.Second("1") != 20001 {
		t.Error(`unexpected value from method Second("1")`)
	}
}

func TestExampleThree(t *testing.T) {
	genjector.Clean()

	err := genjector.Bind(genjector.InSlice(genjector.AsReference[ExampleInterface, *ExampleOne]()))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	err = genjector.Bind(genjector.InSlice(genjector.AsValue[ExampleInterface, ExampleTwo]()))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	err = genjector.Bind(genjector.InSlice(genjector.AsProvider[ExampleInterface](func() (*ExampleOne, error) {
		return &ExampleOne{
			first:  "example3",
			second: 30000,
		}, nil
	})))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	err = genjector.Bind(genjector.InSlice(genjector.AsInstance[ExampleInterface](&ExampleOne{
		first:  "example4",
		second: 40000,
	})))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	instance, err := genjector.Initialize[[]ExampleInterface]()
	if err != nil {
		t.Error("initialization should not cause an error")
	}

	if instance[0].First() != "example" {
		t.Error("unexpected value from method First()")
	}

	if instance[0].Second("1") != 10001 {
		t.Error(`unexpected value from method Second("1")`)
	}

	if instance[1].First() != "example2" {
		t.Error("unexpected value from method First()")
	}

	if instance[1].Second("1") != 20001 {
		t.Error(`unexpected value from method Second("1")`)
	}

	if instance[2].First() != "example3" {
		t.Error("unexpected value from method First()")
	}

	if instance[2].Second("1") != 30001 {
		t.Error(`unexpected value from method Second("1")`)
	}

	if instance[3].First() != "example4" {
		t.Error("unexpected value from method First()")
	}

	if instance[3].Second("1") != 40001 {
		t.Error(`unexpected value from method Second("1")`)
	}
}

func TestExampleFour(t *testing.T) {
	genjector.Clean()

	err := genjector.Bind(genjector.InMap("first", genjector.AsReference[ExampleInterface, *ExampleOne]()))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	err = genjector.Bind(genjector.InMap("second", genjector.AsValue[ExampleInterface, ExampleTwo]()))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	err = genjector.Bind(genjector.InMap("third", genjector.AsProvider[ExampleInterface](func() (*ExampleOne, error) {
		return &ExampleOne{
			first:  "example3",
			second: 30000,
		}, nil
	})))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	err = genjector.Bind(genjector.InMap("fourth", genjector.AsInstance[ExampleInterface](&ExampleOne{
		first:  "example4",
		second: 40000,
	})))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	instance, err := genjector.Initialize[map[string]ExampleInterface]()
	if err != nil {
		t.Error("initialization should not cause an error")
	}

	if instance["first"].First() != "example" {
		t.Error("unexpected value from method First()")
	}

	if instance["first"].Second("1") != 10001 {
		t.Error(`unexpected value from method Second("1")`)
	}

	if instance["second"].First() != "example2" {
		t.Error("unexpected value from method First()")
	}

	if instance["second"].Second("1") != 20001 {
		t.Error(`unexpected value from method Second("1")`)
	}

	if instance["third"].First() != "example3" {
		t.Error("unexpected value from method First()")
	}

	if instance["third"].Second("1") != 30001 {
		t.Error(`unexpected value from method Second("1")`)
	}

	if instance["fourth"].First() != "example4" {
		t.Error("unexpected value from method First()")
	}

	if instance["fourth"].Second("1") != 40001 {
		t.Error(`unexpected value from method Second("1")`)
	}
}

type ExampleThree struct {
	first  string
	second int
	child  ExampleInterface
}

func (e *ExampleThree) Init() {
	e.first = "parent"
	e.second = 100000
	e.child = genjector.MustInitialize[ExampleInterface](genjector.AnnotatedWith("child"))
}

func (e *ExampleThree) First() string {
	return e.first + e.child.First()
}

func (e *ExampleThree) Second(arg string) int {
	value, _ := strconv.ParseInt(arg, 10, 64)
	return e.second + int(value) + e.child.Second(arg)
}

func TestExampleFive(t *testing.T) {
	genjector.Clean()

	err := genjector.Bind(genjector.AsReference[ExampleInterface, *ExampleThree](), genjector.WithAnnotation("parent"))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	err = genjector.Bind(genjector.AsReference[ExampleInterface, *ExampleOne](), genjector.WithAnnotation("child"))
	if err != nil {
		t.Error("binding should not cause an error")
	}

	instance, err := genjector.Initialize[ExampleInterface](genjector.AnnotatedWith("parent"))
	if err != nil {
		t.Error("initialization should not cause an error")
	}

	first := instance.First()
	if first != "parentexample" {
		t.Errorf(`unexpected value from method First(): "%s" != "parentexample"`, first)
	}

	second := instance.Second("1")
	if second != 110002 {
		t.Errorf(`unexpected value from method Second("1"): %d != 110002`, second)
	}
}
