module "pipeline" {
  source = "../infrastructure"
  pipeline_bucket_name = "dev-beeceej-pipelines"
}

resource "aws_iam_user" "pipeline_kickoff" {
  name = "kickoff_blog_post_pipeline"
}

resource "aws_iam_user_policy" "pipeline_kickoff" {
  name = "kickoff-pipeline"
  user = "${aws_iam_user.pipeline_kickoff.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "stateMachine:StartExecution*"
      ],
      "Effect": "Allow",
      "Resource": " arn:aws:states:us-east-1:782123507683:stateMachine:blog-post-pipeline"
    }
  ]
}
EOF
}