module "static_bucket" {
  source = "./modules/static_hosting_s3_bucket"
  bucket_name = "${var.static_bucket_name}"
}
resource "cloudflare_record" "static_CNAME" {
  domain = "beeceej.com"
  name   = "${var.static_bucket_name}"
  value  = "${var.static_bucket_name}.s3-website-us-east-1.amazonaws.com"
  type   = "CNAME"
  proxied = true
}