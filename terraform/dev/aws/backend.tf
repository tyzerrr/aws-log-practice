resource "aws_s3_bucket" "remote_backend"{
  bucket = "${local.project}-remote-backend-${local.env}"

  lifecycle {
    prevent_destroy = true
  }
}

resource "aws_s3_bucket_versioning" "remote_backend" {
  bucket = aws_s3_bucket.remote_backend.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_public_access_block" "remote_backend" {
  bucket = aws_s3_bucket.remote_backend.id
  block_public_acls = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

