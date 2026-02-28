package incident

import "context"

type Metrics struct {
	TotalIncidents   int            `json:"total_incidents"`
	ByState          map[State]int  `json:"by_state"`
	TotalExecutions  int            `json:"total_executions"`
	ExecutionSuccess int            `json:"execution_success"`
	Rollbacks        int            `json:"rollbacks"`
	SuccessRate      float64        `json:"success_rate"`
}

func (s *DynamoStore) GetMetrics(ctx context.Context) (*Metrics, error) {

	items, err := s.ScanAll(ctx)
	if err != nil {
		return nil, err
	}

	metrics := &Metrics{
		ByState: make(map[State]int),
	}

	for _, inc := range items {

		metrics.TotalIncidents++
		metrics.ByState[inc.State]++

		// Count executions
		if inc.State == Executed ||
			inc.State == Verified ||
			inc.State == RolledBack {
			metrics.TotalExecutions++
		}

		// Count success
		if inc.State == Verified {
			metrics.ExecutionSuccess++
		}

		// Count rollbacks
		if inc.State == RolledBack {
			metrics.Rollbacks++
		}
	}

	// Calculate success rate safely
	if metrics.TotalExecutions > 0 {
		metrics.SuccessRate =
			float64(metrics.ExecutionSuccess) /
				float64(metrics.TotalExecutions)
	}

	return metrics, nil
}