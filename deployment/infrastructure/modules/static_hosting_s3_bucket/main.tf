variable "bucket_name" {}

resource "aws_s3_bucket" "static" {
  bucket = "${var.bucket_name}"
  acl    = "public-read"

  website {
    index_document = "index.html"
  }
}