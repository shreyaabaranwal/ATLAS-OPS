resource "aws_iam_role "atlas_ec2_role" {
    name = "${var.project_name}-ec2-role"

    assume_role_policy = jsonencode({
        Version = "2012-10-17"

        Statement = [
            {
                Effect = "Allow"
                Principal = {

                Service = "ec2.amazonaws.com"
            } 
            Action = "sts:AssumeRole"
            }

        ]
    })
}