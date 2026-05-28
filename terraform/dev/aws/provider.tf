provider "aws" {
  region = local.aws_region

  default_tags {
    tags = {
      Env = local.env
      Project = local.project
      ManagedBy = "terraform"
    }
  }
}
