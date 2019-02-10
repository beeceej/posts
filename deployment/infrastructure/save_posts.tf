data aws_iam_policy_document "save_posts" {
  statement {
    actions = ["dynamodb:*"]

    resources = ["${module.posts-dynamodb-table.arn}"]
  }

  statement {
    actions = ["s3:*"]

    resources = ["arn:aws:s3:::static.beeceej.com/*",
      "arn:aws:s3:::beeceej-pipelines/*",
    ]
  }
}

resource "aws_iam_role_policy" "save_posts" {
  name   = "save_posts-policy-attachment"
  role   = "${module.save_posts.lambda_role_name}"
  policy = "${data.aws_iam_policy_document.save_posts.json}"
}

module "save_posts" {
  source        = "./modules/go_lambda"
  function_name = "${local.state_machine_name}-save_posts"
  handler       = "/bin/save_posts"
  file_name     = "../../bin/save_posts.zip"
  memory_size   = "512"
  timeout       = "60"

  environment_vars = {
    "BUCKET_NAME"          = "static.beeceej.com"
    "POSTS_REPO_URL"       = "https://github.com/beeceej/posts"
    "INFLIGHT_BUCKET_NAME" = "beeceej-pipelines"
    "PIPELINE_SUB_PATH"    = "blog-post-pipeline"
    "POSTS_REPO_URI"       = "https://github.com/beeceej/posts"
    "POSTS_TABLE_NAME"     = "${local.table_name}"
  }
}
