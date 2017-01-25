package ginmon

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataStore(t *testing.T) {
	size := 5
	ds := NewDataStore()
	if assert.NotNil(t, ds, "Return of NewDataStore() should not be nil") {
		t.Logf("Should be a dataStore %s", checkMark)
	}

	for j := 0; j < size; j++ {
		for i := 0; i < 100; i++ {
			ds.Add(fmt.Sprintf("%d", i), float64(i))
		}
	}

	l := ds.Get("5")
	if assert.Equal(t, size, len(l), "Return of dataStore#Get() does not work, expect %d but got %d %s",
		size, len(l), ballotX) {
		t.Logf("Return of dataStore#Get() works as expected %d %s",
			len(l), checkMark)
	}

	ds.ResetKey("5")
	l = ds.Get("5")
	if assert.Equal(t, 0, len(l), "Return of dataStore#Get() does not work, expect %d but got %d %s",
		0, len(l), ballotX) {
		t.Logf("Return of dataStore#Get() works as epxected %d %s",
			len(l), checkMark)
	}
}

func TestGenericChannelAspect(t *testing.T) {
	gca := NewGenericChannelAspect("foo")
	if assert.NotNil(t, gca, "Return of NewGenericChannelAspect() should not be nil") {
		t.Logf("Should be a pointer to GenericChannelAspect %s", checkMark)
	}

	l := 100
	for i := 0; i <= l; i++ {
		gca.add(DataChannel{Name: "bar", Value: float64(i)})
	}

	epsilon := 0.01
	gca.calculate()
	gcd := gca.Gcd["bar"]
	if assert.NotNil(t, gcd, "If GenericChannelAspect is not empty and after calling GenericChannelAspect#cacculate() the related GenericChannelData should not be nil") {
		t.Logf("Should be a GenericChannelData %s", checkMark)
	}

	if assert.Equal(t, l+1, gcd.Count, "Count does not work, expect %d but got %d %s",
		l+1, gcd.Count, ballotX) {
		t.Logf("Count works, expected %d %s", gcd.Count, checkMark)
	}
	if assert.Equal(t, 0.0, gcd.Min, "Min does not work, expect %d but got %d %s",
		0.0, gcd.Min, ballotX) {
		t.Logf("Min works, expected %v %s", gcd.Min, checkMark)
	}
	if assert.Equal(t, 100.0, gcd.Max, "Max does not work, expect %d but got %d %s",
		100.0, gcd.Max, ballotX) {
		t.Logf("Max works, expected %v %s", gcd.Max, checkMark)
	}
	if assert.Equal(t, 50.0, gcd.Mean, "Mean does not work, expect %d but got %d %s",
		50.0, gcd.Mean, ballotX) {
		t.Logf("Mean works, expected %v %s", gcd.Mean, checkMark)
	}
	if assert.Equal(t, 90.0, gcd.P90, "P90 does not work, expect %d but got %d %s",
		90.0, gcd.P90, ballotX) {
		t.Logf("P90 works, expected %v %s", gcd.P90, checkMark)
	}
	if assert.Equal(t, 95.0, gcd.P95, "P95 does not work, expect %d but got %d %s",
		95.0, gcd.P95, ballotX) {
		t.Logf("P99 works, expected %v %s", gcd.P99, checkMark)
	}
	if assert.Equal(t, 99.0, gcd.P99, "P99 does not work, expect %d but got %d %s",
		99.0, gcd.P99, ballotX) {
		t.Logf("P99 works, expected %v %s", gcd.P99, checkMark)
	}

	if assert.InEpsilon(t, 29.3, gcd.Stdev, epsilon, "Stdev does not work, expect %d but got %d %s",
		29.3, gcd.Stdev, ballotX) {
		t.Logf("Should be 29.01 %s", checkMark)
	}

}
