package models

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/divergence/pkg/common"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Candle struct {
	Open        float64 `json:"open,string"  validate:"required"`
	High        float64 `json:"high,string"  validate:"required"`
	Low         float64 `json:"low,string"   validate:"required"`
	Close       float64 `json:"close,string" validate:"required"`
	BaseVolume  float64 `json:"baseVolume,string" validate:"required"`
	QuoteVolume float64 `json:"quoteVolume,string" validate:"required"`
	OpenTime    string  `json:"openTime,date" validate:"required"`
	CloseTime   int64   `json:"closeTime,date" validate:"required"`
}

func NewCandle() *Candle {
	candle := &Candle{}

	return candle
}

type Asset struct {
	Date         []time.Time
	Opening      []float64
	Closing      []float64
	High         []float64
	Low          []float64
	Volume       []float64
	VolumeInt    []int64
	Change       []float64
	OpenInterest []float64
}

func (a *Asset) AddCandle(candle *Candle) {
	a.Date = append([]time.Time{time.Unix(common.StringToInt64(candle.OpenTime)/1000, 0)}, a.Date...)
	a.Opening = append([]float64{candle.Open}, a.Opening...)
	a.Closing = append([]float64{candle.Close}, a.Closing...)
	a.High = append([]float64{candle.High}, a.High...)
	a.Low = append([]float64{candle.Low}, a.Low...)
	a.Volume = append([]float64{candle.BaseVolume}, a.Volume...)
	a.VolumeInt = append([]int64{int64(candle.BaseVolume)}, a.VolumeInt...)
	change := (candle.Close - candle.Open) / candle.Open * 100
	a.Change = append([]float64{change}, a.Change...)
}


func PlotCandlestickChart(data []float64, dates []time.Time, market string) {
	p := plot.New()

	p.Title.Text = fmt.Sprintf("Candlestick chart for %s", market)
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

	// Add a legend
	p.Legend.Add("Close", lineClose)

	// Save the plot to a file or display it
	if err := p.Save(6*vg.Inch, 4*vg.Inch, "chart.png"); err != nil {
		log.Fatal(err)
	}
}
