provider "aws" {
  region = "us-east-1"
}

provider "cloudflare" {
  email = "jonesbrianc26@gmail.com"
}

terraform {
  backend "s3" {
    bucket = "prod-beeceej-ops"
    key    = "remote-state/pipeline/terraform.tfstate"
  }
}
