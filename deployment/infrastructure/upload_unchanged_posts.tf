data aws_iam_policy_document "upload_unchanged_posts" {
  statement {
    actions = ["s3:*"]

    resources = [
      "arn:aws:s3:::${var.pipeline_bucket_name}/*",
      "arn:aws:s3:::${var.pipeline_bucket_name}",
      "arn:aws:s3:::${var.static_bucket_name}/*",
      "arn:aws:s3:::${var.static_bucket_name}"
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
  memory_size   = "128"
  timeout       = "60"

  environment_vars = {
    "STATIC_BUCKET_NAME"   = "${var.static_bucket_name}"
    "INFLIGHT_BUCKET_NAME" = "${var.pipeline_bucket_name}"
    "PIPELINE_SUB_PATH"    = "${local.pipeline_sub_path}"
    "POSTS_REPO_URI"       = "${local.posts_repo_uri}"
  }
}
