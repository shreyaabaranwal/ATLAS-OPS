variable "aws_region" {
    description = "AWS deployment region"
    type = string 
    default = "eu-north-1"

}

variable "project_name" {
    description = "Project name"
    type = string 
    default = "atlas-ops"

}

variable "instance_type" {
    description = "EC2 instance type"
    type = string 
    default = "t3.micro"
}