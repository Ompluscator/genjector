# Package Genjector
Reflection-free Run-Time Dependency Injection framework for Go 1.18+

Read [blog](https://hindenbug.io/genjector-reflection-free-run-time-dependency-injection-framework-for-go-1-18-1022d134123d) for better examples!

## The Goal
The purpose of The Genjector package is to provide a Dependency
Injection framework without relying on reflection and depending solely
on Go Generics (provided from Go version 1.18).

It supports many different features:
+ Binding concrete implementations to particular interfaces.
+ Binding implementations as references or values.
+ Binding implementations with Provider methods.
+ Binding implementations with concrete instances.
+ Define Binding as singletons.
+ Define annotations for Binding.
+ Define slices and maps of implementations.
+ ...

## Benchmark
While providing the most of known features of Dependency Injection
frameworks, The Genjector Package also delivers top performance
comparing to current widely used Run-Time DI frameworks, (so, solutions
based on code generators are excluded).

```shell
goos: darwin
goarch: amd64
pkg: github.com/ompluscator/genjector/_benchmark
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
Benchmark
Benchmark/github.com/golobby/container/v3
Benchmark/github.com/golobby/container/v3-8         	 2834061	       409.6 ns/op
Benchmark/github.com/goava/di
Benchmark/github.com/goava/di-8                     	 4568984	       261.9 ns/op
Benchmark/github.com/goioc/di
Benchmark/github.com/goioc/di-8                     	19844284	        60.66 ns/op
Benchmark/go.uber.org/dig
Benchmark/go.uber.org/dig-8                         	  755488	      1497 ns/op
Benchmark/flamingo.me/dingo
Benchmark/flamingo.me/dingo-8                       	 2373394	       503.7 ns/op
Benchmark/github.com/samber/do
Benchmark/github.com/samber/do-8                    	 3585386	       336.0 ns/op
Benchmark/github.com/ompluscator/genjector
Benchmark/github.com/ompluscator/genjector-8        	21460600	        55.71 ns/op
Benchmark/github.com/vardius/gocontainer
Benchmark/github.com/vardius/gocontainer-8          	60947049	        20.25 ns/op
Benchmark/github.com/go-kata/kinit
Benchmark/github.com/go-kata/kinit-8                	  733842	      1451 ns/op
Benchmark/github.com/Fs02/wire
Benchmark/github.com/Fs02/wire-8                    	25099182	        47.43 ns/op
PASS
```

## Examples
Detailed examples can be found inside inner "examples" package.

Some simple code blocks can be found bellow.

### Simple Reference
```go
package example

type ServiceInterface interface {
  String() string
}

type Service struct {
  value string
}

func (s *Service) Init() {
  s.value = "value provided inside the Service"
}

func (s *Service) String() string {
  return s.value
}

err := genjector.Bind(genjector.AsReference[ServiceInterface, *Service]())
if err != nil {
  return err
}

instance, err := genjector.Initialize[ServiceInterface]()
if err != nil {
  return err
}

value := instance.String()
if value != "value provided inside the Service" {
  return err
}
```

### Complex Reference
```go
package example

type ServiceInterface interface {
  String() string
}

type Service struct {
  value string
}

func (s *Service) Init() {
  s.value = "value provided inside the Service"
}

func (s *Service) String() string {
  return s.value
}

err := genjector.Bind(
  genjector.AsReference[ServiceInterface, *Service](),
  genjector.AsSingleton(),
  genjector.WithAnnotation("service")
)
if err != nil {
  return err
}

instance, err := genjector.Initialize[ServiceInterface](
  genjector.WithAnnotation("service"),
)
if err != nil {
  return err
}

value := instance.String()
if value != "value provided inside the Service" {
  return err
}
```
