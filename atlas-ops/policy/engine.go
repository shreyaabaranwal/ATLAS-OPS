package policy

import (
	"context"

	"atlas-ops/infra"
)

type Decision struct {
	Status         string  `json:"status"`
	Confidence     float64 `json:"confidence"`
	AverageCPU     float64 `json:"average_cpu"`
	Slope          float64 `json:"slope"`
	Recommendation string  `json:"recommendation"`
}

func EvaluateTrend(store *infra.MetricsStore, instanceID string) Decision {

	ctx := context.Background()

	metrics, err := store.GetLast5Minutes(ctx, instanceID)
	if err != nil || len(metrics) == 0 {
		return Decision{
			Status: "INSUFFICIENT_DATA",
		}
	}

	avg := calculateAverage(metrics)
	slope := calculateSlope(metrics)
	conf := calculateConfidence(avg, slope, len(metrics))

	if avg > 70 && conf > 0.6 {
		return Decision{
			Status:         "HIGH_CPU",
			Confidence:     conf,
			AverageCPU:     avg,
			Slope:          slope,
			Recommendation: "Scale up instance",
		}
	}

	return Decision{
		Status:     "HEALTHY",
		Confidence: conf,
		AverageCPU: avg,
		Slope:      slope,
	}
}

func calculateAverage(metrics []infra.Metric) float64 {

	if len(metrics) == 0 {
		return 0
	}

	sum := 0.0
	for _, m := range metrics {
		sum += m.CPU
	}

	return sum / float64(len(metrics))
}

func calculateSlope(metrics []infra.Metric) float64 {

	if len(metrics) < 2 {
		return 0
	}

	first := metrics[0].CPU
	last := metrics[len(metrics)-1].CPU

	return last - first
}

func calculateConfidence(avg float64, slope float64, samples int) float64 {

	loadFactor := avg / 100.0

	trendFactor := 0.0
	if slope > 0 {
		trendFactor = 0.3
	}

	durationFactor := float64(samples) / 10.0

	conf := loadFactor*0.5 + trendFactor + durationFactor*0.2

	if conf > 1 {
		return 1
	}

	return conf
}