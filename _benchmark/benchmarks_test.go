package _benchmark_test

import (
	"flamingo.me/dingo"
	"github.com/Fs02/wire"
	"github.com/go-kata/kinit/kinitx"
	"github.com/goava/di"
	di2 "github.com/goioc/di"
	"github.com/golobby/container/v3"
	"github.com/ompluscator/genjector"
	"github.com/samber/do"
	"github.com/vardius/gocontainer"
	"go.uber.org/dig"
	"reflect"
	"testing"
)

type BenchmarkInterface interface {
	Method() string
}

type BenchmarkStruct struct {
	value string
}

func (s *BenchmarkStruct) Method() string {
	return s.value
}

func Benchmark(b *testing.B) {
	b.Run("github.com/golobby/container/v3", func(b *testing.B) {
		var variable BenchmarkInterface

		cont := container.New()
		cont.Transient(func() BenchmarkInterface {
			return &BenchmarkStruct{}
		})
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			cont.Resolve(&variable)
		}

		variable.Method()
	})

	b.Run("github.com/goava/di", func(b *testing.B) {
		var variable BenchmarkInterface

		container, _ := di.New(di.Provide(func() *BenchmarkStruct {
			return &BenchmarkStruct{}
		}, di.As(new(BenchmarkInterface))))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			container.Resolve(&variable)
		}

		variable.Method()
	})

	b.Run("github.com/goioc/di", func(b *testing.B) {
		var variable BenchmarkInterface

		di2.RegisterBean("BenchmarkInterface", reflect.TypeOf(new(BenchmarkStruct)))
		di2.InitializeContainer()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			variable = di2.GetInstance("BenchmarkInterface").(BenchmarkInterface)
		}

		variable.Method()
	})

	b.Run("go.uber.org/dig", func(b *testing.B) {
		var variable BenchmarkInterface

		container := dig.New()
		container.Provide(func() *BenchmarkStruct {
			return &BenchmarkStruct{}
		}, dig.As(new(BenchmarkInterface)))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			container.Invoke(func(result BenchmarkInterface) {
				variable = result
			})
		}

		variable.Method()
	})

	b.Run("flamingo.me/dingo", func(b *testing.B) {
		var variable BenchmarkInterface

		injector, _ := dingo.NewInjector()
		injector.Bind(new(BenchmarkInterface)).To(&BenchmarkStruct{})
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			value, _ := injector.GetInstance(new(BenchmarkInterface))
			variable = value.(BenchmarkInterface)
		}

		variable.Method()
	})

	b.Run("github.com/samber/do", func(b *testing.B) {
		var variable BenchmarkInterface

		injector := do.New()
		do.Provide(injector, func(injector *do.Injector) (BenchmarkInterface, error) {
			return &BenchmarkStruct{}, nil
		})
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			variable = do.MustInvoke[BenchmarkInterface](injector)
		}

		variable.Method()
	})

	b.Run("github.com/ompluscator/genjector", func(b *testing.B) {
		var variable BenchmarkInterface

		genjector.MustBind(genjector.AsReference[BenchmarkInterface, *BenchmarkStruct]())
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			variable = genjector.MustInitialize[BenchmarkInterface]()
		}

		variable = &BenchmarkStruct{}
		variable.Method()
	})

	b.Run("github.com/vardius/gocontainer", func(b *testing.B) {
		var variable BenchmarkInterface

		cont := gocontainer.New()
		cont.Register("BenchmarkInterface", &BenchmarkStruct{})
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			variable = cont.MustGet("BenchmarkInterface").(BenchmarkInterface)
		}

		variable.Method()
	})

	b.Run("github.com/go-kata/kinit", func(b *testing.B) {
		var variable BenchmarkInterface

		kinitx.Provide(func() BenchmarkInterface {
			return &BenchmarkStruct{}
		})
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			kinitx.MustRun(func(benchmarkInterface BenchmarkInterface) {
				variable = benchmarkInterface
			})
		}

		variable.Method()
	})

	b.Run("github.com/Fs02/wire", func(b *testing.B) {
		var variable *BenchmarkStruct

		cont := wire.New()
		cont.Connect(&BenchmarkStruct{})
		cont.Apply()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			cont.Resolve(&variable)
		}

		variable.Method()
	})
}
