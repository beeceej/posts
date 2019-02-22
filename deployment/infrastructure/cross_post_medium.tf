data aws_iam_policy_document "cross_post_medium" {
  statement {
    actions = ["s3:*"]
    effect  = "Allow"

    resources = ["arn:aws:s3:::beeceej-pipelines/*",
      "arn:aws:s3:::static.beeceej.com/*",
    ]
  }
}

resource "aws_iam_role_policy" "cross_post_medium" {
  name   = "cross_post_medium-policy-attachment"
  role   = "${module.convert_posts_to_json.lambda_role_name}"
  policy = "${data.aws_iam_policy_document.convert_posts_to_json.json}"
}

module "cross_post_medium" {
  source        = "./modules/go_lambda"
  function_name = "${local.state_machine_name}-cross_post_medium"
  handler       = "/bin/cross_post_medium"
  file_name     = "../../bin/cross_post_medium.zip"
  memory_size   = "128"
  timeout       = "60"

  environment_vars = {
    "MEDIUM_INTEGRATION_TOKEN" = "${data.aws_ssm_parameter.medium_integration_token.value}"
    "INFLIGHT_BUCKET_NAME"     = "beeceej-pipelines"
    "PIPELINE_SUB_PATH"        = "blog-post-pipeline"
  }
}

data "aws_ssm_parameter" "medium_integration_token" {
  name = "/medium/beeceej.code/integration_token"
}
