# Gin-Gomonitor

[![Build Status](https://travis-ci.org/zalando-techmonkeys/gin-gomonitor.svg?branch=master)](https://travis-ci.org/zalando-techmonkeys/gin-gomonitor)
[![Coverage Status](https://coveralls.io/repos/zalando-techmonkeys/gin-gomonitor/badge.svg?branch=master&service=github)](https://coveralls.io/github/zalando-techmonkeys/gin-gomonitor?branch=master) [![Go Report Card](http://goreportcard.com/badge/zalando-techmonkeys/gin-gomonitor)](http://goreportcard.com/report/zalando-techmonkeys/gin-gomonitor)

Gin-Gomonitor is made specially for [Gin Framework](https://github.com/gin-gonic/gin) users who also want to use [Go-Monitor](https://github.com/mcuadros/go-monitor). It was created by Go developers who needed Gin middleware for exposing metrics with Go-Monitor, which provides a simple and extensible way to build monitorizable long-term execution processes or daemons via HTTP. Gin-Gomonitor supports customized aspects and implements a simple counter aspect within the
package ginmon.

## Project Context and Features

When it comes to choosing a Go framework, there's a lot of confusion about what to use. The scene is very fragmented, and detailed comparisons of different frameworks are still somewhat rare. Meantime, how to handle dependencies and structure projects are big topics in the Go community. 

We've liked using Gin for its speed, accessibility, and usefulness in developing microservice architectures. In creating Gin-Gomonitor, we wanted to take fuller advantage of [Gin](https://github.com/gin-gonic/gin)'s capabilities and help other devs do likewise.

#### How Go-Monitor Is Different from Other Metric Libraries

Go-Monitor is easily extendable, does not need type casts to create JSON, and has useful metrics already defined. It exposes metrics as JSON to a metrics endpoint using a different TCP port.

## Requirements

Gin-Gomonitor uses the following [Go](https://golang.org/) packages as dependencies:

- [Gin](github.com/gin-gonic/gin)
- [Go-Monitor](gopkg.in/mcuadros/go-monitor.v1)

## Installation

Assuming you've installed Go and Gin, run this:

    go get github.com/zalando-techmonkeys/gin-gomonitor

## Usage

[This example](https://github.com/zalando-techmonkeys/gin-gomonitor/blob/master/example/main.go) shows you how to use Gin-Gomonitor. To try it out, use:

      go run example/main.go

### Default Monitor with Custom Aspect

First define your Custom Aspect, so that it will be exposed in a special path /Custom:

```go
type CustomAspect struct {
	CustomValue int
}

func (a *CustomAspect) GetStats() interface{} {
	return a.CustomValue
}

func (a *CustomAspect) Name() string {
	return "Custom"
}

func (a *CustomAspect) InRoot() bool {
	return false
}
```

Next, initialize the CounterAspect defined by Gin-Gomonitor and your own CustomAspect (defined above):

```go
    router := gin.New()

    counterAspect := &ginmon.CounterAspect{0}
    anotherAspect := &CustomAspect{3}
    asps := []aspects.Aspect{counterAspect, anotherAspect}
    router.Use(ginmon.CounterHandler(counterAspect))
```

Finally, register the middleware to expose all metrics on TCP port 9000:

```go
    router.Use(gomonitor.Metrics(9000, asps))
```

#### Testing

The page's counter metric will increment if you hit the page:

    % curl http://localhost:9000/Counter
    {
      "Counter": 0
    }
    % curl http://localhost:8080/
    {"title":"Counter - Hello World - Loook at http://localhost:9000/Counter"}
    % curl http://localhost:8080/
    {"title":"Counter - Hello World - Loook at http://localhost:9000/Counter"}
    % curl http://localhost:9000/Counter
    {
      "Counter": 2
    }

The page custom metric will show three as defined:

    % curl http://localhost:9000/Counter
    {
      "Counter": 0
    }

The regular metrics from go-monitor exposes Go process and build information:

    % curl http://localhost:9000/
    {
      "MemStats": {
        "Alloc": 1285680,
        "TotalAlloc": 1285680,
        "Sys": 4458744,
        "Lookups": 36,
        "Mallocs": 6216,
        "Frees": 0,
        "HeapAlloc": 1285680,
       ...
      "Runtime": {
        "GoVersion": "go1.5.1",
        "GoOs": "linux",
        "GoArch": "amd64",
        "CpuNum": 4,
        "GoroutineNum": 7,
        "Gomaxprocs": 4,
        "CgoCallNum": 1
      },
      "Time": {
        "StartTime": "2016-03-06T17:33:59.440088805+01:00",
        "RequestTime": "2016-03-06T17:48:55.245710327+01:00"
      }
    }

You can also filter by sub-category:

    % curl localhost:9000/Runtime
    {
      "Runtime": {
        "GoVersion": "go1.5.1",
        "GoOs": "linux",
        "GoArch": "amd64",
        "CpuNum": 4,
        "GoroutineNum": 7,
        "Gomaxprocs": 4,
        "CgoCallNum": 1
      }
    }

## Contributing/TODO

We welcome contributions from the community—just submit a pull request. To help you get started, here are some items that we'd love help with:

- Adding more custom metrics
  - time per request: average, p99, ..
  - number of requests: all, path, httpverb
  - &lt;your idea&gt;
- the code base

Please use GitHub issues as your starting point for contributions, new ideas or bug reports.

## Contact

* E-Mail: team-techmonkeys@zalando.de
* IRC on freenode: #zalando-techmonkeys

## Contributors

Thanks to:

- &lt;your name&gt;


## License

See [LICENSE](LICENSE) file.
