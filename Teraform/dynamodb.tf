resource "aws_dynamodb_table" "incidents" {
    name = "${var.project_name}-incidents"
    billing_mode = "PAY_PER_REQUEST"
    hash_key = "IncidentID"


attribute {
    name = "IncidentID"
    type = "S"

}

tags = {
    Project = var.project_name
}

}