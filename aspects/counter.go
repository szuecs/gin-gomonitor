package ginmon

import "github.com/gin-gonic/gin"

// CounterHandler is a Gin middleware function that increments a
// global counter on each request.
func CounterHandler(counter *CounterAspect) gin.HandlerFunc {
	return func(c *gin.Context) {
		counter.Inc()
		c.Next()
	}
}

// CounterAspect stores a counter
type CounterAspect struct {
	Count int
}

func (a *CounterAspect) Inc() {
	a.Count++
}

func (a *CounterAspect) GetStats() interface{} {
	return a.Count
}

func (a *CounterAspect) Name() string {
	return "Counter"
}

func (a *CounterAspect) InRoot() bool {
	return false
}
