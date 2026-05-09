resource "aws_dynamodb_table" "incidents" {
  name         = "atlas-incidents"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "IncidentID"

  attribute {
    name = "IncidentID"
    type = "S"
  }

  tags = {
    Project = "atlas-ops"
  }
}

resource "aws_dynamodb_table" "audit" {
  name         = "atlas-audit"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "incident_id"

  attribute {
    name = "incident_id"
    type = "S"
  }

  tags = {
    Project = "atlas-ops"
  }
}

resource "aws_dynamodb_table" "metrics" {
  name         = "atlas-metrics"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "instance_id"

  attribute {
    name = "instance_id"
    type = "S"
  }

  tags = {
    Project = "atlas-ops"
  }
}