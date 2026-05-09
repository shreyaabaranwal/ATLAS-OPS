<div align="center">

# ⚡ ATLAS-OPS
### Cloud Incident Detection & Automated Remediation System

*Detect infrastructure anomalies, evaluate policy, and execute recovery — with minimal human intervention.*

[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![AWS](https://img.shields.io/badge/AWS-Cloud-FF9900?style=for-the-badge&logo=amazonaws&logoColor=white)](https://aws.amazon.com)
[![DynamoDB](https://img.shields.io/badge/DynamoDB-Storage-4053D6?style=for-the-badge&logo=amazondynamodb&logoColor=white)](https://aws.amazon.com/dynamodb)
[![SQS](https://img.shields.io/badge/SQS-Queue-FF4F8B?style=for-the-badge&logo=amazonsqs&logoColor=white)](https://aws.amazon.com/sqs)
[![Terraform](https://img.shields.io/badge/Terraform-IaC-7B42BC?style=for-the-badge&logo=terraform&logoColor=white)](https://terraform.io)
[![React](https://img.shields.io/badge/React-Dashboard-61DAFB?style=for-the-badge&logo=react&logoColor=black)](https://reactjs.org)

</div>

---

##  What is ATLAS-OPS?

**ATLAS-OPS** is a Go-based cloud incident remediation system built as a backend engineering project to explore SRE automation patterns.

It ingests infrastructure metrics, detects anomalies using a trend-analysis policy engine (EMA + slope + confidence scoring), and routes remediation actions through an SQS-backed async worker pipeline to execute against AWS EC2 resources — with an approval workflow so operators stay informed before critical actions execute.

> This is a functional prototype with core components implemented and Terraform-validated infrastructure. End-to-end production runtime testing is ongoing.

---

##  The Problem

In cloud operations, detection is fast — but response is slow and manual.

| Reality | Cost |
|--------|------|
| Alerts fire; humans get paged | Delayed response, high MTTR |
| Same incidents diagnosed from scratch each time | No automated institutional memory |
| Manual remediation under pressure | Human error, inconsistent outcomes |
| Repetitive low-severity incidents consume on-call bandwidth | Engineer fatigue on fixable problems |

**ATLAS-OPS explores automating the detect → evaluate → act loop** for common EC2 failure patterns, reducing the need for manual intervention on well-understood incidents.

---

##  What's Implemented

| Component | Status | Notes |
|-----------|--------|-------|
|  **Incident Detection** | ✅ Implemented | API ingests metrics; incidents created and classified by severity |
|  **Policy Engine** | ✅ Implemented | EMA + slope + confidence scoring — not just static thresholds |
|  **SQS Async Queue** | ✅ Implemented | Actions enqueued to AWS SQS; decoupled from detection |
|  **Worker Execution Layer** | ✅ Implemented | Goroutine-based workers consume queue and execute EC2 remediations |
|  **DynamoDB Persistence** | ✅ Implemented | Incidents, metrics, and audit logs persisted across all lifecycle states |
|  **Audit Trail** | ✅ Implemented | Every state transition logged via `audit.go` |
|  **EC2 Remediation** | ✅ Implemented | Restart + rollback logic with retry handling via Go SDK v2 |
|  **Approval Workflow** | ✅ Implemented | Operator approval gate before executing critical remediation actions |
|  **Terraform IaC** | ✅ Validated | `init / validate / plan / apply / destroy` tested; all AWS resources provisioned |
|  **React Dashboard** | 🔧 In Progress | Frontend scaffolded; full integration with backend in progress |
|  **CloudWatch Integration** | 🔧 In Progress | Integration implemented; live metric ingestion under validation |

---

## 🏛️ System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    ATLAS-OPS SYSTEM                     │
│                                                         │
│  ┌──────────────┐    ┌───────────────────────────────┐  │
│  │  CloudWatch  │    │        React Dashboard        │  │
│  │  EC2 Metrics │    │   (Incidents · Stats · Logs)  │  │
│  └──────┬───────┘    └───────────────────────────────┘  │
│         │                           ▲                   │
│         ▼                           │                   │
│  ┌──────────────────┐    ┌──────────┴──────────┐        │
│  │ Incident Service │───▶│   DynamoDB Store    │        │
│  │  (api/server.go) │    │  (incident history) │        │
│  └──────┬───────────┘    └─────────────────────┘        │
│         │                                               │
│         ▼                                               │
│  ┌──────────────────────────────────────┐               │
│  │  Policy Engine (policy/engine.go)    │               │
│  │  EMA + slope trend analysis          │               │
│  │  Confidence scoring → action routing │               │
│  └──────┬───────────────────────────────┘               │
│         │                                               │
│         ▼                                               │
│  ┌──────────────────┐                                   │
│  │  Approval Gate   │  ← operator confirms before       │
│  │  (api/server.go) │    critical actions execute       │
│  └──────┬───────────┘                                   │
│         │                                               │
│         ▼                                               │
│  ┌──────────────────┐                                   │
│  │   SQS Queue      │  ← async, decoupled delivery      │
│  │  (queue/sqs.go)  │                                   │
│  └──────┬───────────┘                                   │
│         │                                               │
│         ▼                                               │
│  ┌──────────────────┐    ┌─────────────────────────┐    │
│  │  Worker Layer    │───▶│     EC2 Execution       │    │
│  │(worker/processor)│    │  Restart · Rollback     │    │
│  └──────────────────┘    └─────────────────────────┘    │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

---

##  Project Structure

```
ATLAS-OPS/
│
├── api/
│   ├── server.go              # HTTP server, route registration, approval workflow
│   └── worker.go              # Background worker process entrypoint
│
├── aws/
│   ├── cloudwatch.go          # CloudWatch metrics ingestion
│   └── ec2.go                 # EC2 instance control via AWS SDK Go v2
│
├── cmd/
│   ├── api/
│   │   └── main.go            # API server entry point
│   └── worker/
│       └── main.go            # Worker process entry point
│
├── execution/
│   └── ec2_scaler.go          # EC2 remediation executor with retry + rollback logic
│
├── incident/
│   ├── model.go               # Incident struct and type definitions
│   ├── audit.go               # Audit logging for every lifecycle state transition
│   ├── dynamo_store.go        # DynamoDB read/write operations
│   └── metrics.go             # Incident metrics aggregation
│
├── infra/
│   └── metrics_store.go       # Infrastructure-level metrics persistence
│
├── policy/
│   └── engine.go              # EMA + slope + confidence scoring policy evaluator
│
├── queue/
│   └── sqs.go                 # SQS producer (SendMessage) and consumer logic
│
├── worker/
│   └── processor.go           # Queue consumer, task executor, retry handling
│
├── atlas-ops-dashboard/       # React observability frontend (in progress)
│
├── terraform/
├── go.mod
└── go.sum
```

---

##  Incident Lifecycle

```
1. DETECT
   Metrics ingested via API or CloudWatch
   Threshold + trend analysis → Incident created with severity tag
         │
         ▼
2. PERSIST
   Incident written to DynamoDB (status: OPEN)
   Audit log entry created via audit.go
         │
         ▼
3. EVALUATE
   policy/engine.go computes EMA, slope, and confidence score
   Decision output: RESTART / ROLLBACK / ESCALATE / NO_ACTION
         │
         ▼
4. APPROVE
   For critical actions: operator approval required via API
   Approved → proceeds to queue
   Rejected → incident marked CANCELLED, audit logged
         │
         ▼
5. ENQUEUE
   Approved action published to AWS SQS (message body: incidentID)
   Queue absorbs load; execution stays decoupled from detection
         │
         ▼
6. EXECUTE
   worker/processor.go drains the SQS queue
   Calls execution/ec2_scaler.go + aws/ec2.go
   Retry logic handles transient AWS API failures
   Rollback triggered if execution fails past retry limit
         │
         ▼
7. RESOLVE
   Incident updated in DynamoDB (status: RESOLVED / FAILED)
   Metrics counters updated via incident/metrics.go
   Final audit log entry written
```

---

##  Policy Engine — How Decisions Are Made

Rather than static threshold rules, the policy engine uses **metric trend analysis**:

- **EMA (Exponential Moving Average)** — smooths out noise in incoming metric values to avoid false positives from transient spikes
- **Slope calculation** — detects whether a metric is trending upward, stable, or recovering
- **Confidence scoring** — weights the decision based on signal strength before committing to an action

This means a CPU value of 85% with a rising slope and high confidence triggers remediation; the same value with a flat or declining slope may not — matching how a human SRE would reason about it.

---

##  Distributed Systems Concepts Applied

- **Producer-Consumer Pattern** — Policy engine produces actions; workers consume independently
- **Asynchronous Execution** — SQS decouples detection from remediation; neither layer blocks the other
- **Event-Driven Architecture** — An incident event cascades through detection → policy → approval → queue → execution
- **At-Least-Once Delivery** — SQS visibility timeouts ensure messages survive worker restarts; retry logic in `processor.go` handles redelivery
- **Failure Isolation** — A crashed worker doesn't affect the API server or policy engine
- **Audit Trail** — Every state transition logged via `audit.go` for full incident forensics
- **Separation of Concerns** — Detection, decision, approval, and execution are independent layers

---

##  Tech Stack

### Backend
| Technology | Role |
|------------|------|
| **Go (Golang)** | Core backend — API server, policy engine, worker processes |
| `net/http` | REST API server |
| Goroutines | Concurrent worker execution |
| AWS SDK Go v2 | Direct AWS service integration |

### Cloud — AWS
| Service | Role |
|---------|------|
| **EC2** | Compute target — restart and rollback remediations |
| **SQS** | Async remediation queue — decoupled task delivery |
| **DynamoDB** | Incident, metrics, and audit log persistence |
| **CloudWatch** | Metrics ingestion (integration in progress) |
| **IAM** | Role-based access control for all service identities |

### Infrastructure as Code
| Tool | Role |
|------|------|
| **Terraform** | All AWS resources provisioned declaratively (`init/validate/plan/apply/destroy` validated) |

### Frontend
| Technology | Role |
|------------|------|
| **React** | Observability dashboard (in progress) |

---

##  Getting Started

### Prerequisites

```bash
go 1.21+
aws configure        # AWS CLI with credentials and region configured
terraform 1.5+
node 18+             # For the React dashboard
```

### 1. Clone & Install

```bash
git clone https://github.com/shreyaabaranwal/ATLAS-OPS.git
cd ATLAS-OPS
go mod tidy
```

### 2. Configure Environment

```bash
export AWS_REGION=us-east-1
export DYNAMODB_TABLE=atlas-incidents
export SQS_QUEUE_URL=https://sqs.us-east-1.amazonaws.com/<ACCOUNT_ID>/atlas-queue
```

### 3. Provision Infrastructure

```bash
cd terraform/
terraform init
terraform plan       # Review what will be created
terraform apply      # Provisions SQS, DynamoDB, IAM roles, security groups
```

### 4. Run the API Server

```bash
go run cmd/api/main.go
# Listening on http://localhost:8080
```

### 5. Run the Worker

```bash
# Separate terminal
go run cmd/worker/main.go
# Worker polling SQS for remediation tasks
```

### 6. Start the Dashboard

```bash
cd atlas-ops-dashboard/
npm install && npm start
# Dashboard at http://localhost:3000
```

---

## 🔌 API Reference

### POST `/incident` — Report an Incident

```json
{
  "instance_id": "i-0abc123def456789",
  "severity": "critical",
  "metric": "cpu",
  "value": 95.4
}
```

**Response:**

```json
{
  "incident_id": "INC-20240115-001",
  "status": "OPEN",
  "action": "RESTART_INSTANCE",
  "requires_approval": true
}
```

### POST `/incident/{id}/approve` — Approve Remediation

```json
{ "approved": true }
```

### GET `/incident/{id}` — Fetch Incident State

```json
{
  "incident_id": "INC-20240115-001",
  "instance_id": "i-0abc123def456789",
  "severity": "critical",
  "status": "RESOLVED",
  "action_taken": "RESTART_INSTANCE",
  "resolved_at": "2024-01-15T03:43:12Z"
}
```

### GET `/metrics` — Aggregate Metrics

```json
{
  "total_incidents": 12,
  "resolved": 9,
  "failed": 1,
  "pending_approval": 2
}
```

---

##  Terraform Infrastructure

All AWS resources are provisioned declaratively. Terraform lifecycle (`init / validate / plan / apply / destroy`) has been tested end-to-end.

| Resource | Purpose |
|----------|---------|
| SQS queue | Async task delivery with configurable visibility timeout |
| DynamoDB table | Incident and audit log storage |
| IAM roles | Scoped service identities for API and worker processes |
| IAM instance profiles | EC2-to-AWS service access |
| Security groups | Network access control for provisioned resources |

---

##  Security Considerations

- **IAM roles scoped by service** — API server and worker run under separate IAM roles with only the permissions each needs
- **No hardcoded credentials** — AWS access is via environment variables or IAM instance profiles; no keys in source code
- **`terraform.tfvars` gitignored** — environment-specific values excluded from version control
- **Approval gate for critical actions** — destructive remediations require explicit operator approval before execution

> Note: Security group rules and IAM policies should be reviewed against your specific environment before deploying to any shared or production account.

---

##  Engineering Challenges

- **Trend-based policy design** — moving beyond static thresholds to EMA + slope + confidence scoring to reduce false positives without losing detection sensitivity
- **Decoupled remediation pipeline** — designing the policy, approval, queue, and worker layers to fail and recover independently
- **Approval workflow integration** — threading operator approval into the async pipeline without blocking the detection layer
- **Retry and rollback logic** — handling transient AWS API failures in the executor while avoiding repeated harmful actions
- **Audit consistency** — ensuring every DynamoDB state transition has a corresponding audit log entry

---

##  Future Scope

| Feature | Description |
|---------|-------------|
|  **Dead Letter Queue (DLQ)** | Route exhausted retries to a DLQ for inspection and replay |
|  **Slack / PagerDuty Alerts** | Escalation notifications for incidents that can't be auto-resolved |
|  **Prometheus + Grafana** | Replace custom metrics layer with a standard observability stack |
|  **Kubernetes Integration** | Extend remediation to pod restarts and deployment rollbacks |
|  **Chaos Engineering Tests** | Synthetic failure injection to validate pipeline resilience |
|  **Multi-Region Support** | Cross-region incident routing and remediation |
|  **End-to-End Load Testing** | Validate full incident lifecycle under simulated concurrent load |

---

##  Why ATLAS-OPS?

Most backend projects are a REST API over a database. ATLAS-OPS is an attempt to build something closer to how real SRE systems work — where detection, policy, and execution are separate concerns connected through an async pipeline, and every action is auditable.

The design decisions here (EMA-based policy scoring, approval gates, decoupled worker architecture, full audit trail) reflect patterns used in real incident management platforms, adapted to a scale appropriate for a self-directed backend project.

---

##  Author

Built by **Shreya Baranwal** — Go backend developer & aspiring cloud architect.

Focused on distributed systems, cloud reliability, and infrastructure automation.

[![GitHub](https://img.shields.io/badge/GitHub-shreyaabaranwal-181717?style=flat&logo=github)](https://github.com/shreyaabaranwal)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Shreya%20Baranwal-0A66C2?style=flat&logo=linkedin)](https://www.linkedin.com/in/shreya-baranwal-1103802a5/)
[![Hashnode](https://img.shields.io/badge/Hashnode-Blog-2962FF?style=flat&logo=hashnode)](https://hashnode.com/@shreyabaranwaal)
[![Twitter](https://img.shields.io/badge/Twitter-@Shreyasher786-1DA1F2?style=flat&logo=twitter)](https://twitter.com/Shreyasher786)
[![Email](https://img.shields.io/badge/Email-shreyabaranwal229@gmail.com-D14836?style=flat&logo=gmail)](mailto:shreyabaranwal229@gmail.com)

---

<div align="center">

*Built to learn how real incident systems think — detect, evaluate, approve, execute, audit.*

</div>
