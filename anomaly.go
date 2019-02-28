package anomaly

import (
	"fmt"
	"math"
	"sort"

	"github.com/gonum/stat/distuv"
	"github.com/negbie/stl"
)

func Detect(series []float64, seasonality int, k, a float64) ([]int, []float64, error) {
	if k < 0.01 || a < 0.01 {
		return nil, nil, fmt.Errorf("k and a must be >= 0.01 but k is %f and a is %f", k, a)
	}
	n := len(series)
	_, seasonal, r, err := stl.Decompose(series, seasonality, stl.OuterLoop(1), stl.InnerLoop(2))
	if err != nil {
		return nil, nil, err
	}
	data := make([]float64, n)
	//dataDecomp := make([]float64, n)
	m := quickMedian(series)
	for i := 0; i < n; i++ {
		data[i] = series[i] - seasonal[i] - m
		//dataDecomp[i] = trend[i] + seasonal[i]
	}
	mo := int(math.Round(float64(n) * k))

	rIdx := make([]int, 0, mo)

	distT := &distuv.StudentsT{
		Mu:    0,
		Sigma: 1,
	}

	dataMedian := quickMedian(data)
	dataMean := mean(data)

	tempDataForMad := make([]float64, n)
	for i := 0; i < n; i++ {
		tempDataForMad[i] = math.Abs(data[i] - dataMedian)
	}

	dataStd := quickMedian(tempDataForMad)

	if math.Abs(dataStd) <= 1e-10 {
		return nil, nil, fmt.Errorf("The variance of the series data is zero")
	}

	ares := make([]float64, n)
	for i := 0; i < n; i++ {
		ares[i] = math.Abs(data[i] - dataMedian)
		ares[i] /= dataStd
	}

	dataStd = variance(data)

	aresOrder := orderIndex(ares)
	medianIndex := n / 2
	left := 0
	right := n - 1
	tempMaxIdx := 0
	currentLen := float64(n)
	for i := 1; i < mo+1; i++ {
		npo := float64(n - i + 1)
		nmo := float64(n - i - 1)
		nm := float64(n - i)
		p := 1.0 - a/(2*npo)
		distT.Nu = nmo
		t := distT.Quantile(p)
		tt := t * t
		lam := (t * nm) / math.Sqrt((nmo+tt)*npo)

		if left >= right {
			break
		}
		if currentLen < 1 {
			break
		}

		// remove the largest
		if math.Abs(data[aresOrder[left]]-dataMedian) > math.Abs(data[aresOrder[right]]-dataMedian) {
			tempMaxIdx = aresOrder[left]
			left++
			medianIndex++
		} else {
			tempMaxIdx = aresOrder[right]
			right--
			medianIndex--
		}

		r := math.Abs((data[tempMaxIdx] - dataMedian) / dataStd)
		// recalculate the dataMean and dataStd
		dataStd = math.Sqrt(((currentLen-1)*(dataStd*dataStd+dataMean*dataMean) - data[tempMaxIdx]*data[tempMaxIdx] -
			((currentLen-1)*dataMean-data[tempMaxIdx])*((currentLen-1)*dataMean-data[tempMaxIdx])/
				(currentLen-2)) / (currentLen - 2))
		dataMean = (dataMean*currentLen - data[tempMaxIdx]) / (currentLen - 1)
		dataMedian = data[aresOrder[medianIndex]]
		currentLen--

		if r > a*100*lam {
			rIdx = append(rIdx, tempMaxIdx)
			//fmt.Println(r, lam, series[tempMaxIdx])
		}
	}
	sort.Ints(rIdx)
	return rIdx, r, nil
}
