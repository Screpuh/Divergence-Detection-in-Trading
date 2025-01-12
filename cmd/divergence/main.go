package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/divergence/pkg/common"
	"github.com/divergence/pkg/logger"
	"github.com/divergence/pkg/models"
	"github.com/divergence/pkg/ta/divergence_detection"
)



func main() {
	logger.Info("Starting calculation of divergences on BTC/USDT!")

	// load candles from data/btc-4h.json and convert them to the the candls object in models folder
	candles := loadCandles("./data/btc-4h.json")

	logger.Infof("Loaded %d candles", len(candles.Date))

	// plot the chart for the closing prices of the candles
	models.PlotCandlestickChart(candles.Closing, candles.Date, "BTC/USDT")

	divergence_detection.CalcDivergence(candles.Closing, candles.Date )
}

func loadCandles(location string) models.Asset {
	jsonFile, err := os.Open(location)

	if err != nil {
		fmt.Println(err)
	}
	logger.Infof("Successfully Opened %s", location)
	byteValue, _ := io.ReadAll(jsonFile)

	defer jsonFile.Close()

	var rawData struct {
		List [][]string `json:"list"`
	}
	common.JSONDecode(byteValue, &rawData)

	if err != nil {
		logger.Errorf("Error decoding json: %v", err)
	}

	candles := []*models.Candle{}
	for _, candle := range rawData.List {
		candles = append(candles, &models.Candle{
			Open:        common.StringToFloat64(candle[1]),
			High:        common.StringToFloat64(candle[2]),
			Low:         common.StringToFloat64(candle[3]),
			Close:       common.StringToFloat64(candle[4]),
			BaseVolume:  common.StringToFloat64(candle[5]),
			QuoteVolume: common.StringToFloat64(candle[6]),
			OpenTime:    candle[0],
			CloseTime:   time.Unix(common.StringToInt64(candle[0])/1000, 0).Add(time.Minute * 60).Unix(),
		})
	}

	candleList := models.Asset{}
	for _, candle := range candles {
		candleList.AddCandle(candle)
	}
	
	return candleList
}
