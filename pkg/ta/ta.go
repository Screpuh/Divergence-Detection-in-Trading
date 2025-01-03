package ta

import (
	"math"
	"time"

	"github.com/cinar/indicator"
)

func CalcMovingAverageConvergenceDivergence(candleClose []float64) (macd, signal []float64) {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.Macd(candleClose[maxCandlesNeeded:])
}

func CalcVolumeWeightedAveragePrice(candleClose []float64, volume []float64) []float64 {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	// period is number of hours since beginning of the day
	currentTime := time.Now()
	startOfDay := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	period := int(currentTime.Sub(startOfDay).Hours())

	return indicator.VolumeWeightedAveragePrice(period, candleClose[maxCandlesNeeded:], volume[maxCandlesNeeded:])
}

func CalcAccumulationDistribution(candleClose, candleHigh, candleLow []float64, candleVolume []float64) []float64 {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.AccumulationDistribution(candleHigh[maxCandlesNeeded:], candleLow[maxCandlesNeeded:], candleClose[maxCandlesNeeded:], candleVolume[maxCandlesNeeded:])
}

func CalcBollingerBands(candleClose []float64) ([]float64, []float64, []float64) {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.BollingerBands(candleClose[maxCandlesNeeded:])
}

func CalcActualTrueRange(candleClose, candleHigh, candleLow []float64) ([]float64, []float64) {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.Atr(14, candleHigh[maxCandlesNeeded:], candleLow[maxCandlesNeeded:], candleClose[maxCandlesNeeded:])
}

func CalcAccelarationBands(candleClose, candleHigh, candleLow []float64) ([]float64, []float64, []float64) {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.AccelerationBands(candleHigh[maxCandlesNeeded:], candleLow[maxCandlesNeeded:], candleClose[maxCandlesNeeded:])
}

func CalcWilliamsR(candleClose, candleHigh, candleLow []float64) []float64 {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.WilliamsR(candleLow[maxCandlesNeeded:], candleHigh[maxCandlesNeeded:], candleClose[maxCandlesNeeded:])
}

func CalcAwesomeOscillator(candleClose, candleHigh, candleLow []float64) []float64 {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.AwesomeOscillator(candleLow[maxCandlesNeeded:], candleHigh[maxCandlesNeeded:])
}

func CalcParabolicSar(candleClose, candleHigh, candleLow []float64) ([]float64, []indicator.Trend) {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.ParabolicSar(candleHigh[maxCandlesNeeded:], candleLow[maxCandlesNeeded:], candleClose[maxCandlesNeeded:])
}

func CalcMa(candleClose []float64, length int) []float64 {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - length

	return indicator.Sma(length, candleClose[maxCandlesNeeded:])
}

func CalcEMa(candleClose []float64, length int) []float64 {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - length

	return indicator.Ema(length, candleClose[maxCandlesNeeded:])
}

func CalcAroon(candleClose, candleHigh, candleLow []float64) ([]float64, []float64) {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 50

	return indicator.Aroon(candleHigh[maxCandlesNeeded:], candleLow[maxCandlesNeeded:])
}

func CalcStochasticOscillator(candleClose, candleHigh, candleLow []float64) ([]float64, []float64) {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 20
	return indicator.StochasticOscillator(candleHigh[maxCandlesNeeded:], candleLow[maxCandlesNeeded:], candleClose[maxCandlesNeeded:])
}

func CalcRSI(candleClose []float64) ([]float64, []float64) {
	candleCount := len(candleClose)
	maxCandlesNeeded := candleCount - 40

	return indicator.Rsi(candleClose[maxCandlesNeeded:])
}

// CalcEfficiencyRatio - candleClose needs to be from new to old (newest candle is at index 0)
func CalcEfficiencyRatio(candleClose []float64, n int) float64 {
	candleCount := len(candleClose)
	if (candleCount - n) < 0 {
		return -1
	}

	priceChange := math.Abs(candleClose[0] - candleClose[n])
	sum := 0.0

	for i := n; i > 0; i-- {
		sum = sum + math.Abs(candleClose[i]-candleClose[i+1])
	}

	return priceChange / sum
}

func standardDeviation(num []float64) (float64, float64) {
	var sum, mean, sd float64
	for i, _ := range num {
		sum += num[i]
	}
	mean = sum / float64(len(num))

	for j, _ := range num {
		sd += math.Pow(num[j]-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(num)))

	//get last item from num
	// round nummer

	upper_std := (num[0] - mean) / sd
	//logger.Debug(upper_std)

	return upper_std, sd
}

func CalcChange(current, previous float64) float64 {
	return (current - previous) / previous * 100
}
