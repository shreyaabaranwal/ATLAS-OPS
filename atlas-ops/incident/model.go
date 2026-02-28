package incident

import "time"

type State string

const (
	Detected State = "DETECTED"
	Proposed State = "PROPOSED"
	Approved State = "APPROVED"
	Executed State = "EXECUTED"
	RolledBack State = "ROLLED_BACK" 
	Verified State = "VERIFIED"
)

type Incident struct {
	ID             string
	InstanceID     string
	InstanceType   string
	CPU            float64
	State          State
	Recommendation string
	CreatedAt      time.Time
}