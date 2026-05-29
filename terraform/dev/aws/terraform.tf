terraform {
  required_version = "~> 1.15.0"

  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 6.46"
    }
  }

  backend "s3" {
    bucket = "aws-log-practice-remote-backend-dev"
    key = "terraform/dev/aws/terraform.state"
    region = "ap-northeast-1"
    encrypt = true
    use_lockfile = true
  }
}
