package ginmon

import (
	"time"

	"github.com/gin-gonic/gin"
)

// CounterHandler is a Gin middleware function that increments a
// global counter on each request.
func CounterHandler(ca *CounterAspect) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ca.Inc(ctx)
		ctx.Next()
		ca.IncCodes(ctx)
	}
}

// CounterAspect stores a counter
type CounterAspect struct {
	internalRequestsSum  int
	internalRequests     map[string]int
	internalRequestCodes map[int]int
	RequestsSum          int            `json:"request_sum_per_minute"`
	Requests             map[string]int `json:"requests_per_minute"`
	RequestCodes         map[int]int    `json:"request_codes_per_minute"`
}

// NewCounterAspect returns a new initialized CounterAspect object.
func NewCounterAspect() *CounterAspect {
	ca := &CounterAspect{}
	ca.internalRequestsSum = 0
	ca.internalRequests = make(map[string]int, 0)
	ca.internalRequestCodes = make(map[int]int, 0)
	return ca
}

// StartTimer will call a forever loop in a goroutine to calculate
// metrics for measurements every d ticks. The parameter of this
// function should normally be 1 * time.Minute, if not it will expose
// unintuive JSON keys (requests_per_minute and
// request_sum_per_minute).
func (ca *CounterAspect) StartTimer(d time.Duration) {
	timer := time.Tick(d)
	go func() {
		for {
			<-timer
			ca.reset()
		}
	}()
}

// Inc will increment internal counters that are not exposed. Counters
// will be exposed if you call reset().
func (ca *CounterAspect) Inc(ctx *gin.Context) {
	ca.internalRequestsSum++
	ca.internalRequests[ctx.Request.URL.Path]++
}

// IncCodes will increment internal counters that are not
// exposed. Counters will be exposed if you call reset().
func (ca *CounterAspect) IncCodes(ctx *gin.Context) {
	ca.internalRequestCodes[ctx.Writer.Status()]++
}

// GetStats to fulfill aspects.Aspect interface, it returns the data
// that will be served as JSON.
func (ca *CounterAspect) GetStats() interface{} {
	return *ca
}

// Name to fulfill aspects.Aspect interface, it will return the name
// of the JSON object that will be served.
func (ca *CounterAspect) Name() string {
	return "Counter"
}

// InRoot to fulfill aspects.Aspect interface, it will return where to
// put the JSON object into the monitoring endpoint.
func (ca *CounterAspect) InRoot() bool {
	return false
}

func (ca *CounterAspect) reset() {
	ca.RequestsSum = ca.internalRequestsSum
	ca.Requests = ca.internalRequests
	ca.RequestCodes = ca.internalRequestCodes
	ca.internalRequestsSum = 0
	ca.internalRequests = make(map[string]int, ca.RequestsSum)
	ca.internalRequestCodes = make(map[int]int, len(ca.RequestCodes))
}
