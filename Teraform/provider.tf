terraform {
    required_version = ">= 1.5.0"

    required_proviers {
        aws = {
            source = "harshicorp/aws"
            version = "~> 5.0"
        }
    }
}

provider "aws" {
    region = var.aws_region
}