package ginmon

import "github.com/gin-gonic/gin"

// middleware
func CounterHandler(counter *CounterAspect) gin.HandlerFunc {
	return func(c *gin.Context) {
		counter.Inc()
		c.Next()
	}
}

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
