resource "aws_sqs_queue" incident_queue" {
    name = "${var.project_name}-incident-queue"

    visibility_timeout_seconds = 30
    message_retention_seconds = 86400

    tags = {
        Project = var.project_name
    }
}