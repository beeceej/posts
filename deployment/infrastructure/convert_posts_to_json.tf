data aws_iam_policy_document "convert_posts_to_json" {
  statement {
    actions = ["s3:*"]
    effect  = "Allow"

    resources = [
      "arn:aws:s3:::${var.pipeline_bucket_name}/*",
      "arn:aws:s3:::${var.static_bucket_name}/*",
    ]
  }

  statement {
    effect = "Allow"

    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:DescribeReservedCapacity",
      "dynamodb:DescribeReservedCapacityOfferings",
      "dynamodb:GetRecords",
    ]

    resources = [
      "${module.posts-dynamodb-table.arn}"
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
