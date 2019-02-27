package anomaly

import (
	"math"
	"testing"
)

func TestSeries(t *testing.T) {
	items := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for k := 0; k < 100; k++ {
		quickselect(items, 5)
		if !equal(items[5], 6) {
			t.Errorf("quickselect TestSeries failed - expected 6 got %v\n", items[5])
		}

		quickselect(items, 6)
		if !equal(items[6], 7) {
			t.Errorf("quickselect TestSeries failed - expected 6 got %v\n", items[6])
		}

		quickselect(items, 0)
		if !equal(items[0], 1) {
			t.Errorf("quickselect TestSeries failed - expected 6 got %v\n", items[0])
		}

		quickselect(items, 9)
		if !equal(items[9], 10) {
			t.Errorf("quickselect TestSeries failed - expected 6 got %v\n", items[9])
		}
	}
}

func TestMedian(t *testing.T) {
	itemsEven := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	me := quickMedian(itemsEven)
	if !equal(me, 10.5) {
		t.Errorf("quickselect TestMedian failed - expected 10.5 got %v\n", me)
	}

	itemsOdd := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	mo := quickMedian(itemsOdd)
	if !equal(mo, 10) {
		t.Errorf("quickselect TestMedian failed - expected 10 got %v\n", mo)
	}
}

func equal(a1, a2 float64) bool {
	epsilon := 1e-13
	if math.Abs(a2-a1) > epsilon*math.Abs(a1) {
		return false
	}
	return true
}
