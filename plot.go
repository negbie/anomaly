package anomaly

import (
	"image/color"
	"math"
	"os"
	"time"

	"golang.org/x/image/font/gofont/gomono"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"github.com/golang/freetype/truetype"
)

var defaultFont vg.Font

func init() {
	font, err := truetype.Parse(gomono.TTF)
	if err != nil {
		panic(err)
	}
	vg.AddFont("gomono", font)
	defaultFont, err = vg.MakeFont("gomono", 12)
	if err != nil {
		panic(err)
	}
}

type dateTicks []time.Time

func (t dateTicks) Ticks(min, max float64) []plot.Tick {
	var retVal []plot.Tick
	for i := math.Trunc(min); i <= max; i++ {
		retVal = append(retVal, plot.Tick{Value: i, Label: t[int(i)].String()})
	}
	return retVal
}

type residChart struct {
	plotter.XYs
	draw.LineStyle
}

func (r *residChart) Plot(c draw.Canvas, p *plot.Plot) {
	xmin, xmax, ymin, ymax := r.DataRange()
	p.Y.Min = ymin
	p.Y.Max = ymax
	p.X.Min = xmin
	p.X.Max = xmax

	trX, trY := p.Transforms(&c)
	zero := trY(0)
	lineStyle := r.LineStyle
	for _, xy := range r.XYs {
		x := trX(xy.X)
		y := trY(xy.Y)
		c.StrokeLine2(lineStyle, x, zero, x, y)
	}
}

func (r *residChart) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(1)
	xmax = math.Inf(-1)
	ymin = math.Inf(1)
	ymax = math.Inf(-1)
	for _, xy := range r.XYs {
		xmin = math.Min(xmin, xy.X)
		xmax = math.Max(xmax, xy.X)
		ymin = math.Min(ymin, xy.Y)
		ymax = math.Max(ymax, xy.Y)
	}
	return
}

func (r *residChart) Thumbnail(c *draw.Canvas) {
	y := c.Center().Y
	c.StrokeLine2(r.LineStyle, c.Min.X, y, c.Max.X, y)
}

func newTSPlot(xs []time.Time, ys []float64, seriesName string) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	xys := make(plotter.XYs, len(ys))
	for i := range ys {
		xys[i].X = float64(xs[i].Unix())
		xys[i].Y = ys[i]
	}
	l, err := plotter.NewLine(xys)
	if err != nil {
		return nil, err
	}
	l.LineStyle.Color = color.RGBA{A: 255} // black
	p.Add(l)
	if seriesName != "" {
		p.Legend.Add(seriesName, l)
		p.Legend.TextStyle.Font = defaultFont
	}

	// dieIfErr(plotutil.AddLines(p, seriesName, xys))
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}
	p.Y.Label.TextStyle.Font = defaultFont
	p.X.Label.TextStyle.Font = defaultFont
	p.X.Tick.Label.Font = defaultFont
	p.Y.Tick.Label.Font = defaultFont
	p.Title.Font = defaultFont
	p.Title.Font.Size = 16
	return p, nil
}

func newResidPlot(xs []time.Time, ys []float64, seriesName string) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	xys := make(plotter.XYs, len(ys))
	for i := range ys {
		xys[i].X = float64(xs[i].Unix())
		xys[i].Y = ys[i]
	}
	r := &residChart{XYs: xys, LineStyle: plotter.DefaultLineStyle}
	r.LineStyle.Color = color.RGBA{A: 255}
	p.Add(r)
	p.Legend.Add(seriesName, r)

	p.Legend.TextStyle.Font = defaultFont
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}
	p.Y.Label.TextStyle.Font = defaultFont
	p.X.Label.TextStyle.Font = defaultFont
	p.X.Tick.Label.Font = defaultFont
	p.Y.Tick.Label.Font = defaultFont
	p.Title.Font.Size = 16
	return p, nil
}

func plotDecomposed(xs []time.Time, d, t, s, r []float64) ([][]*plot.Plot, error) {
	plots := make([][]*plot.Plot, 4)
	dp, err := newTSPlot(xs, d, "Series")
	if err != nil {
		return nil, err
	}
	tp, err := newTSPlot(xs, t, "Trend")
	if err != nil {
		return nil, err
	}
	sp, err := newTSPlot(xs, s, "Seasonal")
	if err != nil {
		return nil, err
	}
	rp, err := newTSPlot(xs, r, "Resid")
	if err != nil {
		return nil, err
	}
	plots[0] = []*plot.Plot{dp}
	plots[1] = []*plot.Plot{tp}
	plots[2] = []*plot.Plot{sp}
	plots[3] = []*plot.Plot{rp}
	return plots, nil
}

func writeToPng(a interface{}, title, filename string, width, height vg.Length) error {
	switch at := a.(type) {
	case *plot.Plot:
		at.Title.Text = title
		return at.Save(width*vg.Centimeter, height*vg.Centimeter, filename)
	case [][]*plot.Plot:
		rows := len(at)
		cols := len(at[0])
		t := draw.Tiles{
			Rows: rows,
			Cols: cols,
		}
		img := vgimg.New(width*vg.Centimeter, height*vg.Centimeter)
		dc := draw.New(img)

		if title != "" {
			at[0][0].Title.Text = title
		}

		canvases := plot.Align(at, t, dc)
		for i := 0; i < t.Rows; i++ {
			for j := 0; j < t.Cols; j++ {
				at[i][j].Draw(canvases[i][j])
			}
		}

		w, err := os.Create(filename)
		if err != nil {
			return err
		}

		png := vgimg.PngCanvas{Canvas: img}
		_, err = png.WriteTo(w)
		return err
	}
	panic("Unreachable")
}

func singlePlot(f []float64, d []time.Time, n string) error {
	plt, err := newTSPlot(d, f, n)
	if err != nil {
		return err
	}
	plt.X.Label.Text = "Time"
	plt.Title.Text = n
	return plt.Save(45*vg.Centimeter, 30*vg.Centimeter, "./"+n+".png")
}
