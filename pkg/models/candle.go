package models

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/divergence/pkg/logger"
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

type OpenInterest struct {
	OpenInterest float64 `json:"openInterest,string" validate:"required"`
	OpenTime     string  `json:"openTime,date" validate:"required"`
}

type Currency struct {
	Name   string `json:"name" validate:"required"`
	Market string `json:"market" validate:"required"`
	Base   string `json:"base" validate:"required"`
	Quote  string `json:"quote" validate:"required"`
}

type CandleGreedy struct {
	Close      float64 `json:"close,string" validate:"required"`
	BaseVolume float64 `json:"baseVolume,string" validate:"required"`
	OpenTime   string  `json:"openTime,date" validate:"required"`
}

func NewCandle() *Candle {
	candle := &Candle{}

	return candle
}

func UnmarshallJSON(cBytes []byte) (*Candle, error) {
	var err error
	var cdl Candle

	err = json.Unmarshal(cBytes, &cdl)
	if err != nil {
		return nil, err
	}

	return &cdl, nil
}
func (c *Candle) ToBytes() ([]byte, error) {
	out, err := json.Marshal(c)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return out, nil
}

var periods = map[string]int{"15m": 15, "1h": 60, "4h": 240, "1d": 1440}

func GetValidPeriod(period string) (int, bool) {
	result, ok := periods[period]
	if !ok {
		return result, false
	}

	return result, true
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
