package anomaly

import (
	"math"
	"sort"

	"github.com/chewxy/stl"
	"github.com/gonum/stat/distuv"
)

func Detect(series []float64, seasonality int, k, a float64) ([]int, []float64, error) {
	data := copyslice(series)
	n := len(data)
	decomp := stl.Decompose(data, seasonality, seasonality*2, stl.Additive(), stl.WithRobustIter(15), stl.WithIter(1))
	if decomp.Err != nil {
		return nil, nil, decomp.Err
	}
	copy(data, decomp.Resid)

	mo := int(math.Round(float64(n) * k))
	rIdx := make([]int, 0, mo)
	st := &distuv.StudentsT{}

	for i := 1; i < mo+1; i++ {
		qm := quickMedian(data)
		//fmt.Println("quickMedian", qm)
		ares := aresMedian(data, qm)
		//fmt.Println("bothAres(data)", mad(ares))
		ares = aresMad(ares, mad(data, qm))
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
	return rIdx, decomp.Resid, nil
}
