Atlas Ops is an AI-native Operations Orchestrator that transforms how engineering teams manage cloud infrastructure. By integrating real-time voice interaction, computer vision for dashboard analysis, and a policy-aware reasoning engine, it transitions DevOps from reactive manual toil to "Human-in-the-Loop" autonomous remediation. Built for the Google Gemini Live and Amazon Nova AI hackathons, it demonstrates a production-grade approach to safe, auditable, and multimodal agentic operations.

## The Problem: The "Toil" Gap
Modern cloud environments generate massive telemetry data, yet incident response remains manual and fragmented. High-stakes "On-Call" shifts suffer from:

Context Switching: Engineers bounce between documentation, CLI, and metrics.

Delayed MTTR: Manual root-cause analysis in complex microservices is slow.

Automation Fear: High-risk actions (scaling/rollbacks) lack deterministic safety nets and dry-run validation.

💡 The Solution: Agentic Operations
Atlas Ops acts as a "Senior SRE sitting next to you." It listens to voice commands via Gemini Live, analyzes real-time CloudWatch/Grafana screens via Vision, and executes safe, policy-governed actions across AWS and GCP. It doesn't just "chat"—it manages an Incident State Machine to ensure every action is verified and reversible.

## Core Features
1. Multimodal Interaction (Live Voice + Vision)
Voice-First Interface: Hands-free incident management using Gemini Live for real-time dialogue and interruption-aware command handling.

Visual Reasoning: The agent "sees" dashboards to detect anomalies (e.g., a "jagged" latency spike) that raw logs might miss.

2. Policy-Aware Reasoning Engine
Confidence-Scored Decisions: Every recommendation includes a 0–100% confidence score and a "Reasoning Summary" (e.g., “85% confidence: DB Connection Exhaustion based on RDS metrics”).

Safety Policies: Hard-coded guardrails prevent high-risk actions (e.g., "Never scale production during peak hours without Admin approval").

3. Incident State Machine
Deterministic lifecycle for every event: DETECTED → ANALYZING → PROPOSED → APPROVED → EXECUTED → VERIFIED.

Ensures the agent never loses context during long-running tasks.

4. Reliability & Execution
Dry-Run Mode: Simulates infrastructure changes (API DryRun flags) before actual execution.

Hybrid Execution (API + UI): Primary execution via AWS SDK/Boto3; falls back to Nova Act UI Automation if APIs are throttled or internal consoles are required.

5. Governance & Observability
Immutable Audit Log: Every decision, voice prompt, and execution result is logged for compliance.

Post-Incident Reports: Automated generation of "Executive Summaries" and "Root Cause Analysis" (RCA) documents.

## System Architecture
```
Plaintext
[ USER ] <---(Voice/Vision)---> [ MULTIMODAL GATEWAY ]
                                       |
                                [ ORCHESTRATOR ] <-----> [ STATE MACHINE ]
                                       |              (DynamoDB / Redis)
        _______________________________|_______________________________
       |                               |                              |
[ REASONING LAYER ]           [ EXECUTION LAYER ]           [ OBSERVABILITY ]
 - Gemini 1.5 Pro (Brain)      - AWS SDK (Boto3)             - CloudWatch / Logs
 - Amazon Nova (Agentic)       - Nova Act (UI Auth)          - Telemetry Dashboard
 - Policy Engine               - Approval Workflows          - Audit Trail (S3)
```

## Implementation Details
Google Gemini Live Agent (GCP)
Model: gemini-1.5-flash-8b for low-latency voice and gemini-1.5-pro for complex visual reasoning.

Implementation: Utilizes Multimodal Live API to handle real-time streaming of audio and video frames from the operator's workstation.

Logic: Gemini acts as the primary "Front-End" for natural language understanding and high-level strategy.

Amazon Nova Implementation (AWS)
Model: Amazon Nova Pro & Amazon Nova Lite.

Nova Act: Utilized for advanced UI-based automation where API coverage is incomplete or legacy consoles are used.

Infrastructure: Orchestrated via AWS Step Functions to maintain the Incident State Machine and AWS Lambda for specialized tool-calling.

## Technology Stack
Language: Python 3.11+, TypeScript (Frontend)

AI/ML: Google Gemini API, Amazon Bedrock (Nova), LangGraph (Agent Orchestration)

Cloud: AWS (Lambda, Step Functions, DynamoDB, CloudWatch), GCP (Vertex AI)

Infrastructure as Code: AWS CDK / Terraform

Communication: WebRTC for Live Voice, WebSockets for State Updates

## Setup & Deployment
Local Development
Clone the Repo: git clone https://github.com/user/ops-copilot

Environment Variables: Create a .env with GOOGLE_API_KEY, AWS_ACCESS_KEY_ID, and AWS_SECRET_ACCESS_KEY.

Install Dependencies: pip install -r requirements.txt

Launch Dashboard: npm install && npm start (within /frontend)

Run Agent: python main.py

Cloud Deployment
Backend: Deploy the FastAPI application to AWS App Runner or Google Cloud Run.

State: Provision DynamoDB tables for incident tracking.

Voice: Configure the Multimodal Live API endpoint.

## Demo Walkthrough
Detection: User shares screen showing a Grafana dashboard with rising 5xx errors.

Diagnosis: User asks (Voice): "What's happening?" Agent analyzes the screen and logs, responding with 92% confidence that the DB is throttled.

Proposal: Agent proposes: "I should increase the RDS instance size. This will cost ~$12/day. Proceed?"

Dry-Run: User says: "Do a dry run." Agent simulates the API call and confirms no IAM conflicts.

Execution & Verification: User approves; Agent scales the DB, monitors the latency drop, and confirms: "Systems back to normal. Report generated."

## Safety, Reliability & Governance
Human-in-the-Loop (HITL): No "Write" actions occur without explicit verbal or UI confirmation.

Least Privilege: Agent uses scoped IAM roles with ResourceTag restrictions to prevent accidental deletions.

Reversibility: Every "Action" playbook includes a corresponding "Rollback" path in the state machine.

## Why This Is Different
Most AI chatbots are Passive/Text-Only. Ops-Copilot is:

Proactive: It watches dashboards and alerts the user before they ask.

Multimodal: It understands the spatial layout of a dashboard, not just JSON text.

Stateful: It remembers the incident context across a 20-minute conversation.

Accountable: It provides a deterministic audit trail that satisfies enterprise SRE requirements.

## Future Improvements
Self-Healing Playbooks: Learning from previous human approvals to suggest more accurate "Confidence" scores.

Cost Optimization Mode: Integrating AWS Cost Explorer to suggest cheaper remediation paths.

Multi-Cloud Sync: Native support for Azure Resource Manager.

Hackathon Eligibility Compliance
Google Gemini Live Agent Challenge: Built using Gemini 1.5 Pro and Multimodal Live API features.

Amazon Nova AI Hackathon: Built using Amazon Nova models via Bedrock and utilizing Nova Act for UI-based automation.

Originality: This project was developed specifically for these challenges as a demonstration of next-generation cloud operations.
