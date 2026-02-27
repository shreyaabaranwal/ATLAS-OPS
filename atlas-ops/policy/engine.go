package policy

type Recommendation struct {
	Status         string  `json:"status"`
	CPU            float64 `json:"cpu"`
	InstanceType   string  `json:"instance_type"`
	Recommendation string  `json:"recommendation"`
}

func EvaluatePolicy(cpu float64, instanceType string) Recommendation {

	rec := Recommendation{
		CPU:          cpu,
		InstanceType: instanceType,
	}

	if cpu > 80 {
		rec.Status = "HIGH_CPU"

		switch instanceType {
		case "t2.micro":
			rec.Recommendation = "Upgrade to t3.medium"
		case "t3.medium":
			rec.Recommendation = "Upgrade to t3.large"
		default:
			rec.Recommendation = "Consider vertical scaling"
		}

		return rec
	}

	if cpu > 50 {
		rec.Status = "MODERATE_LOAD"
		rec.Recommendation = "Monitor closely"
		return rec
	}

	rec.Status = "HEALTHY"
	rec.Recommendation = "No action needed"
	return rec
}