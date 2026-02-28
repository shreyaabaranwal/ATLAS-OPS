package incident

import (
	"context"
)

type Metrics struct {
	TotalIncidents     int               `json:"total_incidents"`
	ByState            map[State]int     `json:"by_state"`
	ExecutionSuccess   int               `json:"execution_success"`
	Rollbacks          int               `json:"rollbacks"`
}

func (s *DynamoStore) GetMetrics(ctx context.Context) (*Metrics, error) {

	// Scan entire table (ok for now, optimize later)
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

		if inc.State == Verified {
			metrics.ExecutionSuccess++
		}

		if inc.State == Detected {
			metrics.Rollbacks++
		}
	}

	return metrics, nil
}