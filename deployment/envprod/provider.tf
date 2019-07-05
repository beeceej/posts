provider "aws" {
  region  = "us-east-1"
  profile = "terraformx-assume-role-dev"
   assume_role {
    role_arn     = "arn:aws:iam::448673940787:role/admin"
    session_name = "terraform"
  }
}

terraform {
  backend "s3" {
  profile = "terraformx-assume-role-dev"
    bucket  = "prod-beeceej-ops"
    key     = "remote-state/pipeline/terraform.tfstate"
    region  = "us-east-1"
    role_arn     = "arn:aws:iam::448673940787:role/admin"
  }
}