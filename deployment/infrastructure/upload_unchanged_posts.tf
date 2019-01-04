data aws_iam_policy_document "upload_unchanged_posts" {
  statement {
    actions = ["s3:*"]

    resources = ["arn:aws:s3:::static.beeceej.com/*",
      "arn:aws:s3:::beeceej-pipelines/*",
    ]
  }
}

resource "aws_iam_role_policy" "upload_unchanged_posts" {
  name   = "upload_unchanged_posts-policy-attachment"
  role   = "${module.upload_unchanged_posts.lambda_role_name}"
  policy = "${data.aws_iam_policy_document.upload_unchanged_posts.json}"
}

module "upload_unchanged_posts" {
  source        = "./modules/go_lambda"
  function_name = "${local.state_machine_name}-upload_unchanged_posts"
  handler       = "/bin/upload_unchanged_posts"
  file_name     = "../../bin/upload_unchanged_posts.zip"
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
