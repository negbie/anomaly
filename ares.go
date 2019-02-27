package anomaly

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/stat"
)

func copyslice(input []float64) []float64 {
	i := make([]float64, len(input))
	copy(i, input)
	return i
}

func sortedCopy(input []float64) []float64 {
	i := copyslice(input)
	sort.Float64s(i)
	return i
}

func median(input []float64) float64 {
	i := sortedCopy(input)
	return stat.Quantile(0.5, stat.Empirical, i, nil)
}

func quickMedian(input []float64) float64 {
	i := copyslice(input)
	n := len(i)
	quickselect(i, n/2)
	if (n % 2) == 0 {
		return (i[n/2-1] + i[n/2]) / 2.0
	}
	return i[n/2]
}

func mad(input []float64, m float64) float64 {
	i := make([]float64, len(input))
	for k, v := range input {
		i[k] = math.Abs(v - m)
	}
	return median(i) * 1.4826
}

func aresMad(input []float64, m float64) []float64 {
	i := make([]float64, len(input))
	for k, v := range input {
		i[k] = (v / m)
	}
	return i
}

func aresMedian(input []float64, m float64) []float64 {
	i := make([]float64, len(input))
	for k, v := range input {
		i[k] = math.Abs(v - m)
	}
	return i
}

func maxIdx(input []float64) (max float64, idx int) {
	max = input[0]
	for i := 1; i < len(input); i++ {
		if input[i] > max {
			max = input[i]
			idx = i
		}
	}
	return max, idx
}
