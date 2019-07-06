locals {
  name                 = "kickoff_blog_post_pipeline"
  pipeline_bucket_name = "dev-beeceej-pipelines"
  static_bucket_name   = "posts-dev.beeceej.com"
}

module "pipeline" {
  source               = "../infrastructure"
  pipeline_bucket_name = "${local.pipeline_bucket_name}"
  static_bucket_name   = "${local.static_bucket_name}"
}

resource "aws_iam_user" "pipeline_kickoff" {
  name = "${local.name}"
}

resource "aws_iam_user_policy" "pipeline_kickoff" {
  name = "${local.name}"
  user = "${aws_iam_user.pipeline_kickoff.name}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "states:StartExecution"
      ],
      "Effect": "Allow",
      "Resource": "${module.pipeline.state_machine_id}*"
    }
  ]
}
EOF
}