package ginmon

import (
	"bytes"
	"encoding/gob"
	"sort"
	"sync"
	"time"
)

// DataChannel is the data you pass into the channel. Using Name we
// will put the Value into the right bucket.
type DataChannel struct {
	Name  string
	Value float64
}

type dataStore struct {
	sync.RWMutex
	data map[string][]float64
}

func NewDataStore() dataStore {
	return dataStore{data: make(map[string][]float64)}
}

func (ds dataStore) ResetKey(key string) {
	ds.Lock()
	defer ds.Unlock()
	ds.data[key] = make([]float64, 0)
}

func (ds dataStore) Get(key string) []float64 {
	ds.RLock()
	defer ds.RUnlock()
	return ds.data[key]
}

func (ds dataStore) Add(key string, value float64) {
	ds.data[key] = append(ds.data[key], value)
}

// GenericChannelAspect, exported fields are used to store json
// fields. All fields are measured in nanoseconds.
type GenericChannelAspect struct {
	gcdLock   sync.RWMutex
	name      string
	tempStore dataStore
	Gcd       map[string]GenericChannelData
}

// GenericChannelData
type GenericChannelData struct {
	Count     int       `json:"count"`
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
	gc.tempStore = NewDataStore()
	gc.Gcd = make(map[string]GenericChannelData, 0)
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

// SetupGenericChannelAspect returns an unbuffered channel for type
// DataChannel, such that you can send arbitrary key (string) value
// (float64) pairs to it.
func (gc *GenericChannelAspect) SetupGenericChannelAspect() chan DataChannel {
	lgc := gc // save gc in closure
	ch := make(chan DataChannel)
	go func() {
		for {
			lgc.add(<-ch)
		}
	}()
	return ch
}

// GetStats to fulfill aspects.Aspect interface, it returns a copy of
// the calculated data set that will be served as JSON.
func (gc *GenericChannelAspect) GetStats() interface{} {
	gc.gcdLock.RLock()
	defer gc.gcdLock.RUnlock()

	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)

	err := enc.Encode(gc.Gcd)
	if err != nil {
		return err
	}

	var cpy map[string]GenericChannelData
	err = dec.Decode(&cpy)
	if err != nil {
		return err
	}

	return cpy
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
	gc.tempStore.Lock()
	defer gc.tempStore.Unlock()

	gc.tempStore.Add(dc.Name, dc.Value)
}

func (gc *GenericChannelAspect) calculate() {
	gc.tempStore.Lock()
	defer gc.tempStore.Unlock()
	for name, list := range gc.tempStore.data {
		sortedSlice := list[:]
		gc.tempStore.data[name] = make([]float64, 0)
		l := len(sortedSlice)

		// if tempStore is empty have to set everything to 0 and update timestamp
		if l < 1 {
			gc.gcdLock.Lock()
			gc.Gcd[name] = GenericChannelData{Timestamp: time.Now()}
			gc.gcdLock.Unlock()
			continue
		}

		sort.Float64s(sortedSlice)
		m := mean(sortedSlice, l)

		gc.gcdLock.Lock()
		gc.Gcd[name] = GenericChannelData{
			Timestamp: time.Now(),
			Count:     l,
			Min:       sortedSlice[0],
			Max:       sortedSlice[l-1],
			Mean:      m,
			Stdev:     correctedStdev(sortedSlice, m, l),
			P90:       p90(sortedSlice, l),
			P95:       p95(sortedSlice, l),
			P99:       p99(sortedSlice, l),
		}
		gc.gcdLock.Unlock()
	}
}
