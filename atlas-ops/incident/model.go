package incident

import "time"

type IncidentState string

const (
	Detected        State = "DETECTED"
	Proposed        State = "PROPOSED"
	Approved        State = "APPROVED"
	Executing       State = "EXECUTING"   // distributed lock
	Simulated       State = "SIMULATED"
	ExecutedState   State = "EXECUTED"
	 VerifiedState   State = "VERIFIED"
	FailedPermanent State = "FAILED_PERMANENT"
)

type Incident struct {
	ID             string    `dynamodbav:"id"`
	InstanceID     string    `dynamodbav:"instance_id"`
	InstanceType   string    `dynamodbav:"instance_type"`
	CPU            float64   `dynamodbav:"cpu"`
	State          State     `dynamodbav:"state"`
	Recommendation string    `dynamodbav:"recommendation"`
	CreatedAt      time.Time `dynamodbav:"created_at"`

	// Circuit breaker safety
	ExecutionAttempts int `dynamodbav:"execution_attempts,omitempty"`

	// Verification fields
	CPUAfter         float64    `dynamodbav:"cpu_after,omitempty"`
	VerifiedAt       *time.Time `dynamodbav:"verified_at,omitempty"`
	ExecutionTimeSec int        `dynamodbav:"execution_time_sec,omitempty"`
}