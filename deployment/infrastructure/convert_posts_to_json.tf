data aws_iam_policy_document "convert_posts_to_json" {
  statement {
    actions = ["s3:*"]

    resources = ["arn:aws:s3:::beeceej-pipelines/*",
      "arn:aws:s3:::static.beeceej.com/*",
    ]
  }
}

resource "aws_iam_role_policy" "convert_posts_to_json" {
  name   = "convert_posts_to_json-policy-attachment"
  role   = "${module.convert_posts_to_json.lambda_role_name}"
  policy = "${data.aws_iam_policy_document.convert_posts_to_json.json}"
}

module "convert_posts_to_json" {
  source        = "./modules/go_lambda"
  function_name = "${local.state_machine_name}-convert_posts_to_json"
  handler       = "/bin/convert_posts_to_json"
  file_name     = "../../bin/convert_posts_to_json.zip"
  memory_size   = "512"
  timeout       = "60"

  environment_vars = {
    "BUCKET_NAME"          = "static.beeceej.com"
    "POSTS_REPO_URL"       = "https://github.com/beeceej/posts"
    "INFLIGHT_BUCKET_NAME" = "beeceej-pipelines"
    "PIPELINE_SUB_PATH"    = "blog-post-pipeline"
    "POSTS_REPO_URI"       = "https://github.com/beeceej/posts"
  }
}
