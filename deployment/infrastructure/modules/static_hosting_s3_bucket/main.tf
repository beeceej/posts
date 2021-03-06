variable "bucket_name" {}
variable "bucket_policy_json" {}

resource "aws_s3_bucket" "static" {
  bucket = "${var.bucket_name}"
  acl    = "public-read"
  policy = "${var.bucket_policy_json}"

  website {
    index_document = "index.html"
  }
}