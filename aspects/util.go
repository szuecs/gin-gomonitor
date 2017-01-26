package ginmon

import "math"

func mean(orderedObservations []float64, l int) float64 {
	res := 0.0
	for i := 0; i < l; i++ {
		res += orderedObservations[i]
	}

	return res / float64(l)
}

func p90(orderedObservations []float64, l int) float64 {
	return percentile(orderedObservations, l, 0.9)
}

func p95(orderedObservations []float64, l int) float64 {
	return percentile(orderedObservations, l, 0.95)
}

func p99(orderedObservations []float64, l int) float64 {
	return percentile(orderedObservations, l, 0.99)
}

// percentile with argument p \in (0,1), l is the length of given orderedObservations
// It does a simple apporximation of an ordered list of observations.
// Formula: sortedSlice[0.95*length(sortedSlice)]
func percentile(orderedObservations []float64, l int, p float64) float64 {
	return orderedObservations[int(p*float64(l))]
}

func correctedStdev(observations []float64, mean float64, l int) float64 {
	var omega float64
	for i := 0; i < l; i++ {
		omega += math.Pow(observations[i]-mean, 2)
	}
	stdev := math.Sqrt(1 / (float64(l) - 1) * omega)
	return stdev
}
