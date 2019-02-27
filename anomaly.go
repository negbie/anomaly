package anomaly

import (
	"math"
	"sort"

	"github.com/gonum/stat/distuv"
)

func Detect(series []float64, k, a float64) []int {
	data := copyslice(series)
	n := len(data)
	mo := int(float64(n) * k) // not the best way to round
	rIdx := make([]int, 0, mo)
	st := &distuv.StudentsT{}

	for i := 1; i < mo+1; i++ {
		m := quickMedian(data)
		//fmt.Println("quickMedian", m)
		ares := aresMedian(data, m)
		//fmt.Println("bothAres(data)", mad(ares))
		ares = aresMad(ares, mad(data, m))
		//fmt.Println("aresDivMad(ares, ma)", mad(ares))
		r, idx := maxIdx(ares)
		//fmt.Println(r, data[idx])

		data = append(data[:idx], data[idx+1:]...)

		npo := float64(n - i + 1)
		nmo := float64(n - i - 1)
		nm := float64(n - i)

		p := 1.0 - (a / (2 * npo))

		st.Sigma = 1
		st.Nu = nmo
		t := st.Quantile(p)
		//t := distuv.StudentsT{Mu: 0, Sigma: 1, Nu: nmo, Src: nil}.Quantile(p)

		tt := t * t
		lam := (t * nm) / math.Sqrt((nmo+tt)*npo)
		if r > lam {
			rIdx = append(rIdx, idx)
			//fmt.Println(r, lam)
		}
	}
	sort.Ints(rIdx)
	return rIdx
}
