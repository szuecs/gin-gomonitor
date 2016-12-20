package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-gomonitor"
	"github.com/zalando/gin-gomonitor/aspects"
	"gopkg.in/mcuadros/go-monitor.v1/aspects"
)

func main() {
	requestAspect := ginmon.NewRequestTimeAspect()
	requestAspect.StartTimer(5 * time.Second)

	counterAspect := &ginmon.CounterAspect{0}
	asps := []aspects.Aspect{counterAspect, requestAspect}
	router := gin.New()
	// curl http://localhost:9000/Counter
	router.Use(ginmon.CounterHandler(counterAspect))
	// curl http://localhost:9000/RequestTime
	router.Use(ginmon.RequestTimeHandler(requestAspect))
	// curl http://localhost:9000/
	router.Use(gomonitor.Metrics(9000, asps))
	// last middleware
	router.Use(gin.Recovery())

	// each request to all handlers like below will increment the Counter
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"title": "Counter - Hello World - Loook at http://localhost:9000/Counter"})
	})

	log.Fatal(router.Run(":8080"))
}
