package incident

import "time"

type State string

const (
	Detected   State = "DETECTED"
	Proposed   State = "PROPOSED"
	Approved   State = "APPROVED"
	Simulated  State = "SIMULATED"
	Executed   State = "EXECUTED"
	Verified   State = "VERIFIED"
	RolledBack State = "ROLLED_BACK"
)

type Incident struct {
	ID             string    `dynamodbav:"id"`
	InstanceID     string    `dynamodbav:"instance_id"`
	InstanceType   string    `dynamodbav:"instance_type"`
	CPU            float64   `dynamodbav:"cpu"`
	State          State     `dynamodbav:"state"`
	Recommendation string    `dynamodbav:"recommendation"`
	CreatedAt      time.Time `dynamodbav:"created_at"`
}