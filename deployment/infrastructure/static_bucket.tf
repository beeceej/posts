module "static_bucket" {
  source             = "./modules/static_hosting_s3_bucket"
  bucket_name        = "${var.static_bucket_name}"
  bucket_policy_json = "${data.aws_iam_policy_document.static_bucket.json}"
}
resource "cloudflare_record" "static_CNAME" {
  domain  = "beeceej.com"
  name    = "${var.static_bucket_name}"
  value   = "${var.static_bucket_name}.s3-website-us-east-1.amazonaws.com"
  type    = "CNAME"
  proxied = true
}

data aws_iam_policy_document "static_bucket" {
  statement {
    actions = ["s3:*"]
    principals {
      type        = "AWS"
      identifiers = ["*"]
    }
    resources = [
      "arn:aws:s3:::${var.static_bucket_name}/*",
    ]
  }
}