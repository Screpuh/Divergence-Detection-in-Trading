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

type DivergenceDetection struct {
	Name        string
	Enabled     bool
	RunOnWarmUp bool
}

func New() *DivergenceDetection {
	strat := &DivergenceDetection{}
	strat.SetDefaults()

	return strat
}

func (s *DivergenceDetection) SetDefaults() {
	s.Name = "divergence_detection_strategy"
	s.Enabled = true
	s.RunOnWarmUp = false
}

func (s *DivergenceDetection) RunStrategy() (error) {

	return nil
}

func calcDivergence(candleClose []float64, dates []time.Time, interval int) int {
	if len(candleClose) < 100 {
		return 0
	}

	trend := candleClose[0] - candleClose[len(candleClose)-1]

	var tempCandleClose []float64
	var tempCandleDates []time.Time

	if interval == 15 {
		tempCandleClose = make([]float64, 90)
		tempCandleDates = make([]time.Time, 90)

	} else {
		tempCandleClose = make([]float64, 44)
		tempCandleDates = make([]time.Time, 44)

	}

	// we need to reverse the slice because talib calculates from oldest to newest
	for i, j := 0, len(tempCandleClose)-1; i < j; i, j = i+1, j-1 {
		tempCandleClose[i], tempCandleClose[j] = candleClose[j], candleClose[i]
	}

	// we need to reverse the slice because talib calculates from oldest to newest
	for i, j := 0, len(tempCandleDates)-1; i < j; i, j = i+1, j-1 {
		tempCandleDates[i], tempCandleDates[j] = dates[j], dates[i]
	}

	rsi := talib.Rsi(tempCandleClose, 14)

	// add 2 copies at the end so we can calculate divergence for last 2 candles
	for i := 0; i < 2; i++ {
		rsi = append(rsi, rsi[len(rsi)-1])
		tempCandleDates = append(tempCandleDates, tempCandleDates[len(tempCandleDates)-1])
		tempCandleClose = append(tempCandleClose, tempCandleClose[len(tempCandleClose)-1])
	}

	//rsi_hh := getHigherHighs(rsi, 5, 2)
	//rsi_hl := getHigherLows(rsi, 5, 2)
	//rsi_ll := getLowerLows(rsi, 5, 2)
	//rsi_lh := getLowerHighs(rsi, 5, 2)

	// remove first 14 values because rsi is not calculated for them
	tempCandleDates = tempCandleDates[14:]
	tempCandleClose = tempCandleClose[14:]
	rsi = rsi[14:]

	maxLen := len(rsi) - 2

	order := 6

	plotLocalHighsAndLows(tempCandleClose[:maxLen], tempCandleDates[:maxLen], order)

	plotDivergence2(tempCandleClose[:maxLen], tempCandleDates[:maxLen], "Price", order)

	plotDivergence2(rsi[:maxLen], tempCandleDates[:maxLen], "RSI", order)

	dataPeaks := getPeaks(tempCandleClose, order, 2)
	rsiPeaks := getPeaks(rsi, order, 2)

	logger.Debug("dataPeaks: ", dataPeaks)
	logger.Debug("rsiPeaks: ", rsiPeaks)

	for i := len(dataPeaks["lows"]) - 1; i > len(dataPeaks["lows"])-4; i-- {
		if trend < 0 {
			if dataPeaks["lows"][i] == -1 && rsiPeaks["lows"][i] == 1 {
				// long
				// logger.Debug("Regular divergence: ", tempCandleDates[i])
				return 1
			}

			if dataPeaks["lows"][i] == 1 && rsiPeaks["lows"][i] == -1 {
				// for now i'm not really liking hidden divergence
				// logger.Debug("Hidden divergence: ", tempCandleDates[i])
				return 1
			}
		}

		if trend > 0 {
			if dataPeaks["highs"][i] == 1 && rsiPeaks["highs"][i] == -1 {
				// hidden bearish
				return -1
			}

			if dataPeaks["highs"][i] == 1 && rsiPeaks["highs"][i] == -1 {
				// regular bearish
				return -1
			}
		}
	}

	return 0
}

func Smooth(data []float64, windowSize int) []float64 {
	if windowSize <= 0 {
		return data // No smoothing if the window size is 0 or negative
	}

	smoothed := make([]float64, len(data))

	for i := 0; i < len(data); i++ {
		start := i - windowSize + 1
		if start < 0 {
			start = 0
		}

		sum := 0.0
		count := 0

		for j := start; j <= i; j++ {
			if j < len(data) {
				sum += data[j]
				count++
			}
		}

		smoothed[i] = sum / float64(count)
	}

	return smoothed
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

	logger.Debug("hhIndex: ", hhIndex)
	logger.Debug("lhIndex: ", lhIndex)
	logger.Debug("llIndex: ", llIndex)
	logger.Debug("hlIndex: ", hlIndex)

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
