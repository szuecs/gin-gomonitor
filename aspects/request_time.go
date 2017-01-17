package ginmon

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestTimeAspect, exported fields are used to store json
// fields. All fields are measured in nanoseconds.
type RequestTimeAspect struct {
	lastMinuteRequestTimes []float64
	Min                    float64   `json:"min"`
	Max                    float64   `json:"max"`
	Mean                   float64   `json:"mean"`
	Stdev                  float64   `json:"stdev"`
	P90                    float64   `json:"p90"`
	P95                    float64   `json:"p95"`
	P99                    float64   `json:"p99"`
	Timestamp              time.Time `json:"timestamp"`
}

// NewRequestTimeAspect returns a new initialized RequestTimeAspect
// object.
func NewRequestTimeAspect() *RequestTimeAspect {
	rt := &RequestTimeAspect{}
	rt.lastMinuteRequestTimes = make([]float64, 0)
	rt.Timestamp = time.Now()
	return rt
}

// StartTimer will call a forever loop in a goroutine to calculate
// metrics for measurements every d ticks.
func (rt *RequestTimeAspect) StartTimer(d time.Duration) {
	timer := time.Tick(d)
	go func() {
		for {
			<-timer
			rt.calculate()
		}
	}()
}

// GetStats to fulfill aspects.Aspect interface, it returns the data
// that will be served as JSON.
func (rt *RequestTimeAspect) GetStats() interface{} {
	return rt
}

// Name to fulfill aspects.Aspect interface, it will return the name
// of the JSON object that will be served.
func (rt *RequestTimeAspect) Name() string {
	return "RequestTime"
}

// InRoot to fulfill aspects.Aspect interface, it will return where to
// put the JSON object into the monitoring endpoint.
func (rt *RequestTimeAspect) InRoot() bool {
	return false
}

// RequestTimeHandler is a middleware function to use in Gin
func RequestTimeHandler(rt *RequestTimeAspect) gin.HandlerFunc {
	_rt := rt // save rt in closure
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		took := time.Now().Sub(now)
		_rt.add(float64(took))
	}
}

func (rt *RequestTimeAspect) add(n float64) {
	rt.lastMinuteRequestTimes = append(rt.lastMinuteRequestTimes, n)
}

func (rt *RequestTimeAspect) calculate() {
	sortedSlice := rt.lastMinuteRequestTimes[:]
	rt.lastMinuteRequestTimes = make([]float64, 0)
	l := len(sortedSlice)
	if l <= 1 {
		return
	}
	sort.Float64s(sortedSlice)

	rt.Timestamp = time.Now()
	rt.Min = sortedSlice[0]
	rt.Max = sortedSlice[l-1]
	rt.Mean = mean(sortedSlice, l)
	rt.Stdev = correctedStdev(sortedSlice, rt.Mean, l)
	rt.P90 = p90(sortedSlice, l)
	rt.P95 = p95(sortedSlice, l)
	rt.P99 = p99(sortedSlice, l)
}
