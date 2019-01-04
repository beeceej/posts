locals {
  table_name = "blog-posts"
}

module "posts-dynamodb-table" {
  source        = "./modules/pay_per_request_dynamo_table"
  table_name = "${local.table_name}"
  hash_key     = "ID"
  hash_key_type = "S"
  range_key    = "MD5"
  range_key_type = "S"
}