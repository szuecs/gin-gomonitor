# Gin-Gomonitor

[![Build Status](https://travis-ci.org/zalando/gin-gomonitor.svg?branch=master)](https://travis-ci.org/zalando/gin-gomonitor)
[![Coverage Status](https://coveralls.io/repos/zalando/gin-gomonitor/badge.svg?branch=master&service=github)](https://coveralls.io/github/zalando/gin-gomonitor?branch=master)
[![Go Report Card](https://goreportcard.com/badge/zalando/gin-gomonitor)](https://goreportcard.com/report/zalando/gin-gomonitor)

Gin-Gomonitor is made specially for [Gin Framework](https://github.com/gin-gonic/gin) users who also want to use [Go-Monitor](https://github.com/mcuadros/go-monitor). It was created by Go developers who needed Gin middleware for exposing metrics with Go-Monitor, which provides a simple and extensible way to build monitorizable long-term execution processes or daemons via HTTP. Gin-Gomonitor supports customized aspects and implements a simple counter aspect within the
package ginmon.

## Project Context and Features

When it comes to choosing a Go framework, there's a lot of confusion
about what to use. The scene is very fragmented, and detailed
comparisons of different frameworks are still somewhat rare. Meantime,
how to handle dependencies and structure projects are big topics in
the Go community.

We've liked using Gin for its speed, accessibility, and usefulness in
developing microservice architectures. In creating Gin-Gomonitor, we
wanted to take fuller advantage of
[Gin](https://github.com/gin-gonic/gin)'s capabilities and help other
devs do likewise.

We implemented the following custom
[Aspects](https://github.com/zalando/gin-gomonitor/tree/master/aspects):

CounterAspect implements a counter for request per time.Duration,
counting the sum of all and for each path independent counters.

RequestTimeAspect implements the measurement of request times
including values for min, max, mean, stdev, p90, p95, p99.

GenericChannelAspect implements a generic method to send key value
pairs to this monitoring facility. This aspect will calculate min,
max, mean, stdev, p90, p95, p99 for all values of a key send to the
channel in a given time frame. It can be used without Gin if you like.

See also our [full example](https://github.com/zalando/gin-gomonitor/blob/master/example/main.go).

#### How Go-Monitor Is Different from Other Metric Libraries

Go-Monitor is easily extendable, does not need type casts to create
JSON, and has useful metrics already defined. It exposes metrics as
JSON to a metrics endpoint using a different TCP port.

## Requirements

Gin-Gomonitor uses the following [Go](https://golang.org/) packages as dependencies:

- [Gin](github.com/gin-gonic/gin)
- [Go-Monitor](gopkg.in/mcuadros/go-monitor.v1)

## Installation

Assuming you've installed [Go](https://golang.org/dl) and
[Gin](https://github.com/gin-gonic/gin), run this:

    % go get -u github.com/zalando/gin-gomonitor

## Usage

[This example](https://github.com/zalando/gin-gomonitor/blob/master/example/main.go) shows you how to use Gin-Gomonitor. To try it out, use:

      % go run example/main.go

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

    counterAspect := ginmon.NewCounterAspect()
    counterAspect.StartTimer(1 * time.Minute)

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
        "Counter": {
            "request_sum_per_minute": 0,
            "requests_per_minute": {},
            "request_codes_per_minute": {}
        }
    }
    % for i in {1..20}; do curl localhost:8080/ &>/dev/null ; curl localhost:8080/foo &>/dev/null ; done; sleep 3; curl http://localhost:9000/Counter
    {
        "Counter": {
            "request_sum_per_minute": 40,
            "requests_per_minute": {
                "/": 20,
                "/foo": 20
            },
            "request_codes_per_minute": {
                "200": 20,
                "404": 20
            }
        }
    }

The page custom metric will show three as defined:

    % curl http://localhost:9000/Custom
    {
      "Custom: 3
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

### CounterAspect

CounterAspect measures requests per configured time.Duration.  It
shows requests as sum, per path and per HTTP code, such that you can
monitor increasing user traffic, changing access patterns of user
traffic and http errors.

```go
func main() {
        // initialize CounterAspect and reset every minute
        counterAspect := ginmon.NewCounterAspect()
        counterAspect.StartTimer(1 * time.Minute)

        asps := []aspects.Aspect{counterAspect}
	router := gin.New()
        // register CounterAspect middleware
	// test: curl http://localhost:9000/Counter
        router.Use(ginmon.CounterHandler(counterAspect))

	// start metrics endpoint
	router.Use(gomonitor.Metrics(9000, asps))
	// last middleware
	router.Use(gin.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{ "hello": "world"})
	})

	log.Fatal(router.Run(":8080"))
}
```

The page's counter metric will increment if you hit the page:

    % curl http://localhost:9000/Counter
    {
        "Counter": {
            "request_sum_per_minute": 0,
            "requests_per_minute": {},
            "request_codes_per_minute": {}
        }
    }
    % for i in {1..20}; do curl localhost:8080/ &>/dev/null ; curl localhost:8080/foo &>/dev/null ; done; sleep 3; curl http://localhost:9000/Counter
    {
        "Counter": {
            "request_sum_per_minute": 40,
            "requests_per_minute": {
                "/": 20,
                "/foo": 20
            },
            "request_codes_per_minute": {
                "200": 20,
                "404": 20
            }
        }
    }


### RequestTimeAspect

RequestTimeAspect measures processing time in the middleware
chain. The request will start in the outermost middleware, in this
example it is the RequestTimeHandler. The request will be passed through all
other middleware handlers and at last to your router endpoint. Your
response from your handler will be passed again to all middleware
handlers. RequestTimeAspect will save all measured time.Duration in a
slice and calculate the next metrics with it each timeDuration you
configured with the parameter to StartTimer(d time.Duration). The
metrics endpoint is configured to be
http://localhost:9000/RequestTime, which will be calculated in this
example code every 5 seconds:

```go
func main() {
        // initialize RequestTimeAspect and calculate every 5 seconds
	requestAspect := ginmon.NewRequestTimeAspect()
	requestAspect.StartTimer(5 * time.Second)
	asps := []aspects.Aspect{requestAspect}

	router := gin.New()
        // register RequestTimeAspect middleware
	// test: curl http://localhost:9000/RequestTime
	router.Use(ginmon.RequestTimeHandler(requestAspect))
	// start metrics endpoint
	router.Use(gomonitor.Metrics(9000, asps))
	// last middleware
	router.Use(gin.Recovery())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{ "hello": "world"})
	})

	log.Fatal(router.Run(":8080"))
}
```

RequestTimeAspect will calculate min, max, mean, standard deviation,
P90, P95, P99 of all measured time.Duration for all your endpoints in
this router group. It also creates a time stamp, such that you
know when the calculation happened.

```bash
% for i in {1..20}; do curl localhost:8080/ &>/dev/null; done; sleep 5; curl localhost:9000/RequestTime
{
  "RequestTime": {
    "count": 20,
    "min": 47098,
    "max": 94502,
    "mean": 62199.75,
    "stdev": 13823.2430381624,
    "p90": 91248,
    "p95": 94502,
    "p99": 94502,
    "timestamp": "2017-01-22T19:59:48.164355177+01:00"
  }
}
```

### GenericChannelAspect

GenericChannelAspect enables you to send arbitrary ginmon.DataChannel
data through a channel to gin-gomonitor, which will calculate min,
max, mean, standard deviation, P90, P95, P99 grouped by
ginmon.DataChannel.Name for every every configured time.Duration. The
metrics endpoint is configured to be http://localhost:9000/generic,
which will be calculated in this example code every 3 seconds:

```go
func main() {
        // initialize GenericChannelAspect and calculate every 3 seconds
	genericAspect := ginmon.NewGenericChannelAspect("generic")
	genericAspect.StartTimer(3 * time.Second)
	genericCH := genericAspect.SetupGenericChannelAspect()
	asps := []aspects.Aspect{genericAspect}

	router := gin.New()
        // register GenericChannelAspect middleware
	// test: curl http://localhost:9000/generic
	// start metrics endpoint
	router.Use(gomonitor.Metrics(9000, asps))
	// catch panics as last middleware
	router.Use(gin.Recovery())

        // send a lot of data concurrently to the monitoring data channel
	i := 0
	go func() {
		for {
			i++
			genericCH <- ginmon.DataChannel{Name: "foo", Value: float64(i)}
		}
	}()
	j := 0
	go func() {
		for {
			j++
			genericCH <- ginmon.DataChannel{Name: "bar", Value: float64(j % 5)}
		}
	}()

	router.Run(":8080"))
}
```

```bash
% curl http://localhost:9000/generic
{
  "generic": {
    "bar": {
      "count": 2110190,
      "min": 0,
      "max": 4,
      "mean": 2,
      "stdev": 1.4142138974647371,
      "p90": 4,
      "p95": 4,
      "p99": 4,
      "timestamp": "2017-01-24T14:40:20.970407737+01:00"
    },
    "foo": {
      "count": 2445672,
      "min": 5.128943e+06,
      "max": 7.574614e+06,
      "mean": 6.3517785e+06,
      "stdev": 706004.8381128459,
      "p90": 7.330047e+06,
      "p95": 7.452331e+06,
      "p99": 7.550158e+06,
      "timestamp": "2017-01-24T14:40:21.299909533+01:00"
    }
  }
}
```

## Contributing/TODO

We welcome contributions from the community—just submit a pull
request. To help you get started, here are some items that we'd love
help with:

- Adding more custom metrics
  - add more tests
  - time per request: path, httpverb
  - number of requests: httpverb
  - review and maybe refactor lock usage in generic_channel.go
  - reduce goroutine usage: We could use one goroutine for all myAspect.StartTimer()
  - add logging and enable user to choose logging, see [Dave Cheney's post](https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern)
  - &lt;your idea&gt;
- the code base

Please use GitHub issues as your starting point for contributions, new ideas or bug reports.

## Contact

* E-Mail: [MAINTAINERS](MAINTAINERS)

## Contributors

Thanks to:

- &lt;your name&gt;


## License

See [LICENSE](LICENSE) file.
