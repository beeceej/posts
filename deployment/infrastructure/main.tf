terraform {
  backend "s3" {
    bucket = "beeceej-ops"
    region = "us-east-1"
  }
}

locals {
  state_machine_name = "blog-post-pipeline"
}

data "template_file" "state_machine_definition" {
  template = "${file("${path.module}/state.json")}"

  vars {
    convert_posts_to_json  = "${module.convert_posts_to_json.lambda_arn}"
    upload_unchanged_posts = "${module.upload_unchanged_posts.lambda_arn}"
    save_posts = "${module.save_posts.lambda_arn}"
  }
}

data aws_iam_policy_document "state_machine_assume_role" {
  statement {
    principals {
      identifiers = [
        "states.us-east-1.amazonaws.com",
      ]

      type = "Service"
    }

    actions = [
      "sts:AssumeRole",
    ]
  }
}

resource "aws_iam_role" "state_machine_role" {
  name               = "${local.state_machine_name}"
  assume_role_policy = "${data.aws_iam_policy_document.state_machine_assume_role.json}"
}

data aws_iam_policy_document "state_machine_policy" {
  statement {
    actions = [
      "lambda:InvokeFunction",
    ]

    resources = [
      "*",
    ]
  }
}

resource "aws_iam_policy" "invoke_lambda_policy" {
  name   = "invoke-lambda-policy"
  policy = "${data.aws_iam_policy_document.state_machine_policy.json}"
}

resource "aws_iam_policy_attachment" "policy_attach" {
  name       = "lambda-invoke-policy-attachment"
  roles      = ["${aws_iam_role.state_machine_role.name}"]
  policy_arn = "${aws_iam_policy.invoke_lambda_policy.arn}"
}

resource "aws_sfn_state_machine" "state_machine" {
  name     = "${local.state_machine_name}"
  role_arn = "${aws_iam_role.state_machine_role.arn}"

  definition = "${data.template_file.state_machine_definition.rendered}"
}

resource "aws_s3_bucket" "inflight_bucket" {
  bucket = "beeceej-pipelines"
}

output "state_machine_id" {
  value = "${aws_sfn_state_machine.state_machine.id}"
}

output "state_machine_name" {
  value = "${aws_sfn_state_machine.state_machine.name}"
}
