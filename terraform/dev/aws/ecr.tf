resource "aws_ecr_repository" "ecr" {
  name                 = "${local.project}-${local.env}-ecr-repository"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "AES256"
  }

  tags = {
    Name = "${local.project}-${local.env}-ecr-repository"
  }
}

resource "aws_ecr_lifecycle_policy" "ecr-lifecycle-policy" {
  repository = aws_ecr_repository.ecr.name
  policy     = file("./policy/ecr_lifecycle_policy.json")
}
