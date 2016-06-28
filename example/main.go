package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-gomonitor"
	"github.com/zalando/gin-gomonitor/aspects"
	"gopkg.in/mcuadros/go-monitor.v1/aspects"
)

func main() {
	counterAspect := &ginmon.CounterAspect{0}
	asps := []aspects.Aspect{counterAspect}
	router := gin.New()
	// curl http://localhost:9000/Counter
	router.Use(ginmon.CounterHandler(counterAspect))
	// curl http://localhost:9000/
	router.Use(gomonitor.Metrics(9000, asps))
	// last middleware
	router.Use(gin.Recovery())

	// each request to all handlers like below will increment the Counter
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"title": "Counter - Hello World - Loook at http://localhost:9000/Counter"})
	})

	//..
	router.Run(":8080")
}
