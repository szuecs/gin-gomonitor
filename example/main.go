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

	counterAspect := ginmon.NewCounterAspect()
	counterAspect.StartTimer(3 * time.Second)

	genericAspect := ginmon.NewGenericChannelAspect("generic")
	genericAspect.StartTimer(3 * time.Second)
	genericCH := genericAspect.SetupGenericChannelAspect()

	asps := []aspects.Aspect{counterAspect, requestAspect, genericAspect}

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
		ctx.JSON(http.StatusOK, gin.H{
			"Counter": map[string]string{
				"msg": "Request Counter - Loook at http://localhost:9000/Counter",
				"cmd": "curl http://localhost:9000/Counter ; for i in {1..20}; do curl localhost:8080/ &>/dev/null ; curl localhost:8080/foo &>/dev/null ; done; sleep 3; curl http://localhost:9000/Counter"},
			"RequestTime": map[string]string{
				"msg": "RequestTime is registered at http://localhost:9000/RequestTime and will return data after 5 seconds.",
				"cmd": "for j in {0..100}; do for i in {1..20}; do curl localhost:8080/ ; done; sleep 0.5; curl localhost:9000/RequestTime ; done"},
			"GenericChannelAspect": map[string]string{
				"msg": "Generic Aspect can process arbitrary map[string]float64 data - Loook at http://localhost:9000/generic",
				"cmd": "curl http://localhost:9000/generic ; for i in {1..20}; do curl localhost:8080/generic &>/dev/null ; done; sleep 3; curl http://localhost:9000/generic"}})
	})

	router.GET("/generic", func(ctx *gin.Context) {
		for i := 0; i < 100; i++ {
			genericCH <- ginmon.DataChannel{Name: "foo", Value: float64(i % 2)}
			genericCH <- ginmon.DataChannel{Name: "bar", Value: float64(i % 5)}
		}
		ctx.JSON(http.StatusOK, gin.H{
			"GenericChannelAspect": map[string]string{
				"msg": "Generic Aspect can process arbitrary map[string]float64 data - Loook at http://localhost:9000/generic",
				"cmd": "curl http://localhost:9000/generic ; for i in {1..20}; do curl localhost:8080/generic &>/dev/null ; done; sleep 3; curl http://localhost:9000/generic"}})
	})

	log.Fatal(router.Run(":8080"))
}
