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

func mean(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	m := 0.0
	for i, x := range xs {
		m += (x - m) / float64(i+1)
	}
	return m
}

func variance(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	} else if len(xs) <= 1 {
		return 0
	}
	mean, M2 := 0.0, 0.0
	for n, x := range xs {
		delta := x - mean
		mean += delta / float64(n+1)
		M2 += delta * (x - mean)
	}
	return M2 / float64(len(xs)-1)
}

type by struct {
	Indices []int
	Values  []float64
}

func (b by) Len() int           { return len(b.Values) }
func (b by) Less(i, j int) bool { return b.Values[i] < b.Values[j] }
func (b by) Swap(i, j int) {
	b.Indices[i], b.Indices[j] = b.Indices[j], b.Indices[i]
	b.Values[i], b.Values[j] = b.Values[j], b.Values[i]
}

func orderIndex(input []float64) []int {
	i := copyslice(input)
	out := make([]int, len(i))

	for k := range i {
		out[k] = k
	}
	sort.Sort(by{Indices: out, Values: i})
	return out
}
