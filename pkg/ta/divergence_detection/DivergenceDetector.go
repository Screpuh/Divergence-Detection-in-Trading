package divergence_detection

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/divergence/pkg/logger"
	"github.com/markcheno/go-talib"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func CalcDivergence(candleClose []float64, dates []time.Time) int {
	var tempCandleClose []float64
	var tempCandleDates []time.Time
	
	tempCandleClose = make([]float64, len(candleClose))
	tempCandleDates = make([]time.Time, len(candleClose))

	// we want to move the candleClose and dates to a variable where we can specify the length
	// lets select the first 50 candles for smaller sample set
	tempCandleClose = candleClose[0:80]
	tempCandleDates = dates[0:80]

	rsi := talib.Rsi(tempCandleClose, 14)

	order := 4

	// remove the first 14 elements from the array, because we don't have RSI values for them
	plotLocalHighsAndLows(tempCandleClose[14:], tempCandleDates[14:], order)
	
	plotDivergence2(tempCandleClose[14:], tempCandleDates[14:], "trend_lines_price", order)

	plotDivergence2(rsi[14:], tempCandleDates[14:], "trend_lines_rsi", order)

	dataPeaks := getPeaks(tempCandleClose, order, 2)
	rsiPeaks := getPeaks(rsi, order, 2)

	divergences := []string{}

	for i := 0; i < len(dataPeaks["lows"]); i++ {
		if dataPeaks["lows"][i] == -1 && rsiPeaks["lows"][i] == 1 {
			// long
			//logger.Debugf("Regular bullish divergence: %v", tempCandleDates[i])
			divergences = append(divergences, fmt.Sprintf("Regular bullish divergence: %v", tempCandleDates[i]))
		}

		if dataPeaks["lows"][i] == 1 && rsiPeaks["lows"][i] == -1 {
			//logger.Debug("Hidden bullis divergence: ", tempCandleDates[i])
			divergences = append(divergences, fmt.Sprintf("Hidden bullish divergence: %v", tempCandleDates[i]))
		}

		if dataPeaks["highs"][i] == -1 && rsiPeaks["highs"][i] == 1 {
			// hidden bearish
			//logger.Debug("Hidden bearish divergence: ", tempCandleDates[i])
			divergences = append(divergences, fmt.Sprintf("Hidden bearish divergence: %v", tempCandleDates[i]))
		}

		if dataPeaks["highs"][i] == 1 && rsiPeaks["highs"][i] == -1 {
			// regular bearish
			//logger.Debug("Regular bearish divergence: ", tempCandleDates[i])
			divergences = append(divergences, fmt.Sprintf("Regular bearish divergence: %v", tempCandleDates[i]))
		}
	}

	logger.Debugf("Divergences: %v", divergences)

	return 0
}

func plotDivergence2(data []float64, dates []time.Time, title string, order int) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Price ($)"

	// Plot the close prices
	ptsClose := make(plotter.XYs, len(data))
	for i := range data {
		ptsClose[i].X = float64(i)
		ptsClose[i].Y = data[i]
	}
	lineClose, err := plotter.NewLine(ptsClose)
	if err != nil {
		log.Fatal(err)
	}
	lineClose.Color = color.RGBA{0, 0, 0, 255}
	p.Add(lineClose)

	// Define colors for different divergence points
	colors := []color.Color{
		color.RGBA{0, 0, 0, 255},     // Close
		color.RGBA{255, 0, 0, 255},   // Higher Highs
		color.RGBA{0, 255, 0, 255},   // Higher Lows
		color.RGBA{0, 0, 255, 255},   // Lower Lows
		color.RGBA{255, 0, 255, 255}, // Lower Highs
	}

	hh := getHigherHighs(data, order, 2)
	hl := getHigherLows(data, order, 2)
	ll := getLowerLows(data, order, 2)
	lh := getLowerHighs(data, order, 2)

	// Plot different divergence points
	plotDivergence(p, hh, dates, data, colors[1])
	plotDivergence(p, hl, dates, data, colors[2])
	plotDivergence(p, ll, dates, data, colors[3])
	plotDivergence(p, lh, dates, data, colors[4])

	// Add a legend
	p.Legend.Add("Close", lineClose)

	// Save the plot to a file or display it
	if err := p.Save(6*vg.Inch, 4*vg.Inch, fmt.Sprintf("%s.png", title)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Plot saved as divergence.png")
}

func plotRsi(data []float64, dates []time.Time) {
	p := plot.New()

	p.Title.Text = "RSI"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "RSI"

	rsi := talib.Rsi(data, 14)

	ptsRsi := make(plotter.XYs, len(rsi))
	for i := range rsi {
		ptsRsi[i].X = float64(i)
		ptsRsi[i].Y = rsi[i]
	}

	lineRsi, err := plotter.NewLine(ptsRsi)
	if err != nil {
		log.Fatal(err)
	}
	lineRsi.Color = color.RGBA{0, 0, 0, 255}
	p.Add(lineRsi)

	// Save the plot to a file or display it
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "rsi.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Plot saved as rsi.png")
}

func plotLocalHighsAndLows(data []float64, dates []time.Time, order int) {
	maxIdx := boolRelExtrema(data, order, func(a, b float64) bool { return a > b })
	minIdx := boolRelExtrema(data, order, func(a, b float64) bool { return a < b })

	p := plot.New()

	p.Title.Text = "Maxima and Minima Points"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Price ($)"

	// Plot the close prices
	ptsClose := make(plotter.XYs, len(data))
	for i := range data {
		ptsClose[i].X = float64(i)
		ptsClose[i].Y = data[i]
	}

	lineClose, err := plotter.NewLine(ptsClose)
	if err != nil {
		log.Fatal(err)
	}
	lineClose.Color = color.RGBA{0, 0, 0, 255}
	p.Add(lineClose)

	// Define colors for maxima and minima points
	maximaColor := color.RGBA{255, 0, 0, 255} // Red
	minimaColor := color.RGBA{0, 0, 255, 255} // Blue

	// Plot maxima points
	plotMaxima(p, maxIdx, dates, data, maximaColor)

	// Plot minima points
	plotMinima(p, minIdx, dates, data, minimaColor)

	// Add a legend
	p.Legend.Add("Close", lineClose)

	// Save the plot to a file or display it
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "maxima_minima.png"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Plot saved as maxima_minima.png")
}

func plotDivergence(p *plot.Plot, indices [][]int, dates []time.Time, close []float64, c color.Color) {
	for _, group := range indices {
		pts := make(plotter.XYs, len(group))
		for i, idx := range group {
			pts[i].X = float64(idx)
			pts[i].Y = close[idx]
		}
		line, err := plotter.NewLine(pts)
		if err != nil {
			log.Fatal(err)
		}
		line.Color = c
		p.Add(line)
	}
}

func plotMaxima(p *plot.Plot, indices []int, dates []time.Time, close []float64, c color.Color) {
	pts := make(plotter.XYs, len(indices))
	for i, idx := range indices {
		pts[i].X = float64(idx)
		pts[i].Y = close[idx]
	}
	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	s.GlyphStyle.Color = c
	s.GlyphStyle.Radius = vg.Points(5)
	p.Add(s)
}

func plotMinima(p *plot.Plot, indices []int, dates []time.Time, close []float64, c color.Color) {
	pts := make(plotter.XYs, len(indices))
	for i, idx := range indices {
		pts[i].X = float64(idx)
		pts[i].Y = close[idx]
	}
	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	s.GlyphStyle.Color = c
	s.GlyphStyle.Radius = vg.Points(5)
	p.Add(s)
}

func getPeaks(data []float64, order, K int) map[string][]float64 {
	hhIndex := getHHIndex(data, order, K)
	lhIndex := getLHIndex(data, order, K)
	llIndex := getLLIndex(data, order, K)
	hlIndex := getHLIndex(data, order, K)

	dataWithPeaks := make(map[string][]float64)

	//create dataWithPeaks["highs"] and dataWithPeaks["lows"] with NaN values length of data
	dataWithPeaks["highs"] = make([]float64, len(data))
	dataWithPeaks["lows"] = make([]float64, len(data))

	for _, idx := range hhIndex {
		dataWithPeaks["highs"][idx] = 1
	}
	for _, idx := range lhIndex {
		dataWithPeaks["highs"][idx] = -1
	}

	for _, idx := range llIndex {
		dataWithPeaks["lows"][idx] = -1
	}
	for _, idx := range hlIndex {
		dataWithPeaks["lows"][idx] = 1
	}

	return dataWithPeaks
}

func getHHIndex(data []float64, order, K int) []int {
	extrema := getHigherHighs(data, order, K)
	logger.Debug("higher high extrema: ", extrema)
	var idx []int
	for _, i := range extrema {
		if i[len(i)-1]+order < len(data) {
			idx = append(idx, i[len(i)-1]+order)
		}
	}
	logger.Debug("higher high idx: ", idx)
	return idx
}

func getLHIndex(data []float64, order, K int) []int {
	extrema := getLowerHighs(data, order, K)
	logger.Debug("lower high extrema: ", extrema)
	var idx []int
	for _, i := range extrema {
		if i[len(i)-1]+order < len(data) {
			idx = append(idx, i[len(i)-1]+order)
		}
	}
	logger.Debug("lower high idx: ", idx)
	return idx
}

func getLLIndex(data []float64, order, K int) []int {
	extrema := getLowerLows(data, order, K)

	logger.Debug("lower low extrema: ", extrema)
	var idx []int
	for _, i := range extrema {
		if i[len(i)-1]+order < len(data) {
			idx = append(idx, i[len(i)-1]+order)
		}
	}
	logger.Debug("lower low idx: ", idx)
	return idx
}

func getHLIndex(data []float64, order, K int) []int {
	extrema := getHigherLows(data, order, K)
	var idx []int
	logger.Debug("higher low extrema: ", extrema)
	for _, i := range extrema {
		if i[len(i)-1]+order < len(data) {
			idx = append(idx, i[len(i)-1]+order)
		}
	}
	logger.Debug("higher low idx: ", idx)
	return idx
}

func getHigherHighs(data []float64, order, K int) [][]int {
	comparator := func(a, b float64) bool {
		return a < b
	}

	highIdx := boolRelExtrema(data, order, comparator)
	highs := []float64{}

	for _, i := range highIdx {
		highs = append(highs, data[i])
	}

	var result [][]int
	var current []int

	for i, _ := range highIdx {
		if i == 0 {
			current = append(current, highIdx[i])
			continue
		}
		if highs[i] > highs[i-1] {
			if len(current) == 0 {
				current = append(current, highIdx[i-1])
			}
			current = append(current, highIdx[i])

			if len(current) == K {
				result = append(result, current)
				current = nil
			}
		} else {
			current = nil
		}
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}

func getLowerHighs(data []float64, order, K int) [][]int {
	comparator := func(a, b float64) bool {
		return a < b
	}

	highIdx := boolRelExtrema(data, order, comparator)
	highs := []float64{}

	for _, i := range highIdx {
		highs = append(highs, data[i])
	}

	var result [][]int
	var current []int

	for i, _ := range highIdx {
		if i == 0 {
			current = append(current, highIdx[i])
			continue
		}
		if highs[i] < highs[i-1] {
			if len(current) == 0 {
				current = append(current, highIdx[i-1])
			}
			current = append(current, highIdx[i])

			if len(current) == K {
				result = append(result, current)
				current = nil
			}
		} else {
			current = nil
		}
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}

func getLowerLows(data []float64, order, K int) [][]int {
	comparator := func(a, b float64) bool {
		return a > b
	}

	lowIdx := boolRelExtrema(data, order, comparator)
	lows := []float64{}

	for _, i := range lowIdx {
		lows = append(lows, data[i])
	}

	var result [][]int
	var current []int

	for i, _ := range lowIdx {
		if i == 0 {
			current = append(current, lowIdx[i])
			continue
		}
		if lows[i] < lows[i-1] {
			if len(current) == 0 {
				current = append(current, lowIdx[i-1])
			}
			current = append(current, lowIdx[i])

			if len(current) == K {
				result = append(result, current)
				current = nil
			}
		} else {
			current = nil
		}
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}

func getHigherLows(data []float64, order, K int) [][]int {
	comparator := func(a, b float64) bool {
		return a > b
	}

	lowIdx := boolRelExtrema(data, order, comparator)
	lows := []float64{}

	for _, i := range lowIdx {
		lows = append(lows, data[i])
	}

	var result [][]int
	var current []int

	for i, _ := range lowIdx {
		if i == 0 {
			current = append(current, lowIdx[i])
			continue
		}
		if lows[i] > lows[i-1] {
			if len(current) == 0 {
				current = append(current, lowIdx[i-1])
			}
			current = append(current, lowIdx[i])

			if len(current) == K {
				result = append(result, current)
				current = nil
			}
		} else {
			current = nil
		}
	}

	if len(current) > 0 {
		result = append(result, current)
	}

	return result
}

func boolRelExtrema(data []float64, order int, comparator func(float64, float64) bool) []int {
	if order < 1 {
		panic("Order must be an int >= 1")
	}

	extrema := make([]bool, len(data))

	for i := order; i < len(data)-order; i++ {
		isExtrema := true

		for j := i - order; j <= i+order; j++ {
			if j != i && comparator(data[i], data[j]) {
				isExtrema = false
				break
			}
		}

		extrema[i] = isExtrema
	}

	return extremaToIndices(extrema)
}

func containsTrue(arr []bool) bool {
	for _, v := range arr {
		if v {
			return true
		}
	}
	return false
}

func extremaToIndices(extrema []bool) []int {
	indices := []int{}
	for i, isExtrema := range extrema {
		if isExtrema {
			indices = append(indices, i)
		}
	}
	return indices
}
