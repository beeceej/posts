provider "aws" {
  region = "us-east-1"
}

terraform {
  backend "s3" {
    bucket = "dev-beeceej-ops"
    key    = "remote-state/pipeline/terraform.tfstate"
  }
}
