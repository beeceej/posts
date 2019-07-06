provider "aws" {
  region = "us-east-1"
}

provider "cloudflare" {
  email = "jonesbrianc26@gmail.com"
  token = "${var.cloudflare_token}"
}

terraform {
  backend "s3" {
    bucket = "dev-beeceej-ops"
    key    = "remote-state/pipeline/terraform.tfstate"
  }
}
