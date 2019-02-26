package anomaly

import (
	"math/rand"
	"time"
)

// Quickselect selects the kth smallest item from array
func Quickselect(array []float64, k int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	quickselect(array, k, 0, len(array), r)
}

// quickselect is the main recursive quicksort routine
func quickselect(array []float64, k int, first int, last int, r *rand.Rand) {
	if first == last-1 {
		return
	}

	pivot := partitionItems(array, first, last, r)
	if k < pivot {
		quickselect(array, k, first, pivot, r)
	} else if k > pivot {
		quickselect(array, k, pivot+1, last, r)
	} else {
		return
	}
}

// partitionItems partitions the items according to a randomly chosen pivot
func partitionItems(array []float64, first, last int, r *rand.Rand) int {
	pivot := first + r.Intn(last-first-1)
	array[first], array[pivot] = array[pivot], array[first]

	j := first + 1
	for i := first + 1; i < last; i++ {
		if array[i] < array[first] {
			array[j], array[i] = array[i], array[j]
			j++
		}
	}
	array[j-1], array[first] = array[first], array[j-1]
	return j - 1
}
