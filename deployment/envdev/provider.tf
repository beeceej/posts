provider "aws" {
  region  = "us-east-1"
  # profile = "terraformx"
  # assume_role {
  #   role_arn     = "arn:aws:iam::782123507683:role/admin"
  #   session_name = "terraform"
  # }
}

terraform {
  backend "s3" {
    # profile  = "terraformx"
    bucket   = "dev-beeceej-ops"
    key      = "remote-state/pipeline/terraform.tfstate"
    region   = "us-east-1"
    # role_arn = "arn:aws:iam::782123507683:role/admin"
  }
}

