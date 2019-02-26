package anomaly

import (
	"math/rand"
	"testing"
)

var items []float64

func init() {
	items = randFloats(1, 1000, 10000)
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func BenchmarkMedian(b *testing.B) {
	data := items
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = median(data)
	}
}

func BenchmarkQuickMedian(b *testing.B) {
	data := items
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = quickMedian(data)
	}
}
