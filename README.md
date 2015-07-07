# gomonitor
[Gin](https://github.com/gin-gonic/gin) middleware for exposing metrics with
[go-monitor](https://github.com/mcuadros/go-monitor). It supports
custom aspects and implements a simple counter aspect in the package ginmon.

[![Build Status](https://travis-ci.org/zalando-techmonkeys/gin-gomonitor.svg?branch=master)](https://travis-ci.org/zalando-techmonkeys/gin-gomonitor)
[![Coverage Status](https://coveralls.io/repos/zalando-techmonkeys/gin-gomonitor/badge.svg?branch=master&service=github)](https://coveralls.io/github/zalando-techmonkeys/gin-gomonitor?branch=master)

## Usage
### Example

```go
package main

import (
	"github.com/gin-gonic/gin"
        "github.com/zalando-techmonkeys/gin-gomonitor"
	"github.com/zalando-techmonkeys/gin-gomonitor/aspects"
	"gopkg.in/mcuadros/go-monitor.v1"
	"gopkg.in/mcuadros/go-monitor.v1/aspects"
)

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

func main() {
    router := gin.New()

    counterAspect := &ginmon.CounterAspect{0}
    anotherAspect := &CustomAspect{3}
    asps := []aspects.Aspect{counterAspect, anotherAspect}

    // curl http://localhost:9000/Counter
    router.Use(ginmon.CounterHandler(counterAspect))
    // curl http://localhost:9000/
    // curl http://localhost:9000/Custom
    router.Use(gomonitor.Metrics(9000, asps))

    // last middleware
    router.Use(gin.Recovery())

    // each request to all handlers like below will increment the Counter
    router.GET("/", func(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"title": "Counter - Hello World"})
    })

    //..
    router.Run(":8080")
}
```

## Development
* Issues: Just create issues on github
* Enhancements/Bugfixes: Pull requests are welcome
* get in contact: mailto:team-techmonkeys@zalando.de
* see [MAINTAINERS](https://github.com/zalando-techmonkeys/gin-gomonitor/blob/master/MAINTAINERS)
file.

## License
see [LICENSE](https://github.com/zalando-techmonkeys/gin-gomonitor/blob/master/LICENSE) file.
