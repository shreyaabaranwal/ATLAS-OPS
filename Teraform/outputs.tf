output "ec2_public_ip" {
    value = aws_instance.atlas_server.public_ip
}

output "sqs_url" {
    value = aws_sqs_queue.incident_queue.url
}

output "dynamodb_table_name" {
    value = aws_dynamodb_table.incidents.name
}