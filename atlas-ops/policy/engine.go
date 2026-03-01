package policy

import (
	"context"
	"math"

	"atlas-ops/infra"
)

type Decision struct {
	Status         string  `json:"status"`
	Confidence     float64 `json:"confidence"`
	AverageCPU     float64 `json:"average_cpu"`
	Slope          float64 `json:"slope"`
	Recommendation string  `json:"recommendation,omitempty"`
}

func EvaluateTrend(store *infra.MetricsStore, instanceID string) Decision {

	ctx := context.Background()

	metrics, err := store.GetLast5Minutes(ctx, instanceID)
	if err != nil || len(metrics) == 0 {
		return Decision{
			Status: "INSUFFICIENT_DATA",
		}
	}

	alpha := 0.3

	ema := calculateEMA(metrics, alpha)
	slope := calculateSlope(metrics)
	current := metrics[len(metrics)-1].CPU
	spike := detectSpike(current, ema)

	conf := calculateConfidence(ema, slope, spike, len(metrics))

	if ema > 75 && conf > 0.7 {
		return Decision{
			Status:         "HIGH_CPU",
			Confidence:     conf,
			AverageCPU:     ema,
			Slope:          slope,
			Recommendation: "Scale up instance",
		}
	}

	return Decision{
		Status:     "HEALTHY",
		Confidence: conf,
		AverageCPU: ema,
		Slope:      slope,
	}
}

// ---------------- EMA ----------------

func calculateEMA(metrics []infra.Metric, alpha float64) float64 {

	if len(metrics) == 0 {
		return 0
	}

	ema := metrics[0].CPU

	for i := 1; i < len(metrics); i++ {
		ema = alpha*metrics[i].CPU + (1-alpha)*ema
	}

	return ema
}

// ---------------- NORMALIZED SLOPE ----------------

func calculateSlope(metrics []infra.Metric) float64 {

	if len(metrics) < 2 {
		return 0
	}

	first := metrics[0].CPU
	last := metrics[len(metrics)-1].CPU

	rawSlope := last - first

	// Normalize slope to prevent over-influence
	return rawSlope / float64(len(metrics))
}

// ---------------- SPIKE DETECTION ----------------

func detectSpike(current float64, ema float64) float64 {

	diff := current - ema

	if diff <= 0 {
		return 0
	}

	// Normalize spike between 0–1
	return math.Min(diff/100.0, 1)
}

// ---------------- CONFIDENCE ENGINE ----------------

func calculateConfidence(
	ema float64,
	slope float64,
	spike float64,
	sampleCount int,
) float64 {

	loadScore := ema / 100.0

	trendScore := 0.0
	if slope > 0 {
		trendScore = math.Min(slope/10.0, 1)
	}

	stabilityScore := math.Min(float64(sampleCount)/10.0, 1)

	confidence :=
		loadScore*0.5 +
			trendScore*0.2 +
			spike*0.2 +
			stabilityScore*0.1

	return math.Min(confidence, 1)
}