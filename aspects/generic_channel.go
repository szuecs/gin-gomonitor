package ginmon

import (
	"sort"
	"time"
)

// DataChannel is the data you pass into the channel. Using Name we
// will put the Value into the right bucket.
type DataChannel struct {
	Name  string
	Value float64
}

// GenericChannelAspect, exported fields are used to store json
// fields. All fields are measured in nanoseconds.
type GenericChannelAspect struct {
	name      string
	tempStore map[string][]float64
	Gcd       map[string]*GenericChannelData
}

// GenericChannelData
type GenericChannelData struct {
	Min       float64   `json:"min"`
	Max       float64   `json:"max"`
	Mean      float64   `json:"mean"`
	Stdev     float64   `json:"stdev"`
	P90       float64   `json:"p90"`
	P95       float64   `json:"p95"`
	P99       float64   `json:"p99"`
	Timestamp time.Time `json:"timestamp"`
}

// NewGenericChannelAspect returns a new initialized GenericChannelAspect
// object.
func NewGenericChannelAspect(name string) *GenericChannelAspect {
	gc := &GenericChannelAspect{name: name}
	gc.tempStore = make(map[string][]float64, 0)
	gc.Gcd = make(map[string]*GenericChannelData, 0)
	return gc
}

// StartTimer will call a forever loop in a goroutine to calculate
// metrics for measurements every d ticks.
func (gc *GenericChannelAspect) StartTimer(d time.Duration) {
	timer := time.Tick(d)
	go func() {
		for {
			<-timer
			gc.calculate()
		}
	}()
}

// TODO: change name and signature
func (gc *GenericChannelAspect) SetupGenericChannelAspect() chan DataChannel {
	_gc := gc // save gc in closure
	_ch := make(chan DataChannel, 1)
	go func(gc *GenericChannelAspect, ch chan DataChannel) {
		for {
			gc.add(<-ch)
		}
	}(_gc, _ch)
	return _ch
}

// GetStats to fulfill aspects.Aspect interface, it returns the data
// that will be served as JSON.
func (gc *GenericChannelAspect) GetStats() interface{} {
	return gc.Gcd
}

// Name to fulfill aspects.Aspect interface, it will return the name
// of the JSON object that will be served.
func (gc *GenericChannelAspect) Name() string {
	return gc.name
}

// InRoot to fulfill aspects.Aspect interface, it will return where to
// put the JSON object into the monitoring endpoint.
func (gc *GenericChannelAspect) InRoot() bool {
	return false
}

func (gc *GenericChannelAspect) add(dc DataChannel) {
	gc.tempStore[dc.Name] = append(gc.tempStore[dc.Name], dc.Value)
}

func (gc *GenericChannelAspect) calculate() {
	for name, list := range gc.tempStore {
		sortedSlice := list[:]
		gc.tempStore[name] = make([]float64, 0)
		l := len(sortedSlice)

		// if tempStore is empty have to set everything to 0 and update timestamp
		if l < 1 {
			gc.Gcd[name] = &GenericChannelData{Timestamp: time.Now()}
			continue
		}

		sort.Float64s(sortedSlice)
		m := mean(sortedSlice, l)

		gc.Gcd[name] = &GenericChannelData{
			Timestamp: time.Now(),
			Min:       sortedSlice[0],
			Max:       sortedSlice[l-1],
			Mean:      m,
			Stdev:     correctedStdev(sortedSlice, m, l),
			P90:       p90(sortedSlice, l),
			P95:       p95(sortedSlice, l),
			P99:       p99(sortedSlice, l),
		}
	}
}
