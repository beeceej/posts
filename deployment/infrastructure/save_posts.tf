data aws_iam_policy_document "save_posts" {
  statement {
    actions = ["dynamodb:*"]

    resources = ["${module.posts-dynamodb-table.arn}"]
  }

  statement {
    actions = ["s3:*"]

    resources = [
      "arn:aws:s3:::${var.static_bucket_name}/*",
      "arn:aws:s3:::${var.pipeline_bucket_name}/*",
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
  memory_size   = "128"
  timeout       = "60"

  environment_vars = {
    "STATIC_BUCKET_NAME"   = "${var.static_bucket_name}"
    "INFLIGHT_BUCKET_NAME" = "${var.pipeline_bucket_name}"
    "PIPELINE_SUB_PATH"    = "${local.pipeline_sub_path}"
    "POSTS_REPO_URI"       = "${var.posts_repo_uri}"
    "POSTS_TABLE_NAME"     = "${local.table_name}"
  }
}
