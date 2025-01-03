package logger

import (
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
	"time"
)

type Point struct {
	X float64
	Y float64
}

func PrintChart(data []float64, highs, lows []Point, symbol string) {
	var xv []time.Time
	var yv []float64
	for x, y := range data {
		xv = append(xv, time.Now().AddDate(0, 0, x))
		yv = append(yv, y)
	}

	min, max := GetExtremes(yv)

	priceSeries := chart.TimeSeries{
		Name: "SPY",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: xv,
		YValues: yv,
	}

	var xvbo []float64
	var yvbo []float64
	var xvso []float64
	var yvso []float64

	for _, point := range highs {
		xvbo = append(xvbo, float64(time.Now().AddDate(0, 0, int(point.X)).UnixNano()))
		yvbo = append(yvbo, point.Y)
	}

	for _, point := range lows {
		xvso = append(xvso, float64(time.Now().AddDate(0, 0, int(point.X)).UnixNano()))
		yvso = append(yvso, point.Y)
	}

	buyOrderSeries := chart.ContinuousSeries{
		Style: chart.Style{
			Show:        true,
			StrokeWidth: chart.Disabled,
			DotWidth:    5,
			DotColor:    drawing.Color{R: 0, G: 255, B: 0, A: 255},
		},
		XValues: xvbo,
		YValues: yvbo,
	}

	sellOrderSeries := chart.ContinuousSeries{
		Style: chart.Style{
			Show:        true,
			StrokeWidth: chart.Disabled,
			DotWidth:    5,
			DotColor:    drawing.Color{R: 255, G: 0, B: 0, A: 255},
		},
		XValues: xvso,
		YValues: yvso,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Max: max,
				Min: min,
			},
		},
		Series: []chart.Series{
			priceSeries,
			buyOrderSeries,
			sellOrderSeries,
		},
	}

	f, _ := os.Create("test1.PNG")
	defer f.Close()
	err := graph.Render(chart.PNG, f)

	if err != nil {
		Error(err)
	}
}

func GetExtremes(data []float64) (float64, float64) {
	var min float64
	var max float64

	for _, point := range data {
		if min == 0 {
			min = point
			max = point
		} else {
			if point < min {
				min = point
			}

			if point > max {
				max = point
			}
		}
	}

	return min, max
}
