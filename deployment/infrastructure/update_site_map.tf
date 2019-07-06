data aws_iam_policy_document "update_site_map" {
  statement {
    actions = ["s3:*"]

    resources = [
      "arn:aws:s3:::${var.static_bucket_name}.com/*",
      "arn:aws:s3:::${var.static_bucket_name}.com",
      "arn:aws:s3:::beeceej-pipelines/*",
      "arn:aws:s3:::beeceej-pipelines"
    ]
  }
}

resource "aws_iam_role_policy" "update_site_map" {
  name   = "update_site_map-policy-attachment"
  role   = "${module.update_site_map.lambda_role_name}"
  policy = "${data.aws_iam_policy_document.update_site_map.json}"
}

module "update_site_map" {
  source        = "./modules/go_lambda"
  function_name = "${local.state_machine_name}-update_site_map"
  handler       = "/bin/update_site_map"
  file_name     = "../../bin/update_site_map.zip"
  memory_size   = "512"
  timeout       = "60"

  environment_vars = {
    "STATIC_BUCKET_NAME"   = "${var.static_bucket_name}"
    "INFLIGHT_BUCKET_NAME" = "${var.pipeline_bucket_name}"
    "PIPELINE_SUB_PATH"    = "${local.pipeline_sub_path}"
    "POSTS_REPO_URI"       = "${local.posts_repo_uri}"
  }
}
