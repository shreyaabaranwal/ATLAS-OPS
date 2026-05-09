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

resource "aws_iam_role_policy_attachment" "cloudwatch" {
    role = aws_iam_role.atlas_ec2_role.name
    policy_arn = "arn:aws:iam::aws:policy/CloudWatchFullAccess"

}

resource "aws_iam_role_policy_attachment" "sqs" {
    role = aws_iam_role.atlas_ec2_role.name
    policy_arn = "arn:aws:iam::aws:policy/AmazonSQSFullAccess"

}

resource "aws_iam_role_policy_attachment" "dynamodb" {
    role = aws_iam_role.atlas_ec2_role.name
    policy_arn = "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"

}

resource "aws_iam_role_policy_attachment" "ec2" {
    role = aws_iam_role.atlas_ec2_role.name
    policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"

}


resource "aws_iam_instance_profile "atlas_profile" {
    name = "${var.project_name}-instance-profile"
    role = aws_iam_role.atlas_ec2_role.name
}