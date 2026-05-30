resource "aws_lb" "alb" {
  name                       = "${local.project}-${local.env}-alb"
  internal                   = false
  enable_deletion_protection = false
  load_balancer_type         = "application"
  subnets                    = [for subnet in aws_subnet.public_subnets : subnet.id]
  security_groups            = [aws_security_group.alb_sg.id]
  access_logs {
    bucket  = aws_s3_bucket.alb_log_bucket.id
    enabled = true
  }
}

# ALB Log bucket
resource "aws_s3_bucket" "alb_log_bucket" {
  bucket = "${local.project}-${local.env}-alb-log-bucket"

  tags = {
    Name = "${local.project}-${local.env}-alb-log-bucket"
  }
}

resource "aws_s3_bucket_public_access_block" "alb_log_bucket_public_access_block" {
  bucket                  = aws_s3_bucket.alb_log_bucket.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_server_side_encryption_configuration" "alb_log_bucket_encryption" {
  bucket = aws_s3_bucket.alb_log_bucket.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_lifecycle_configuration" "alb_log_bucket_lifecycle" {
  bucket = aws_s3_bucket.alb_log_bucket.id
  rule {
    id     = "${local.project}-${local.env}-alb-log-bucket-lifecycle"
    status = "Enabled"

    transition {
      storage_class = "STANDARD_IA"
      days          = 30
    }

    expiration {
      days = 90
    }
  }
}

# ALB bucket iam policy
data "aws_iam_policy_document" "alb_log_bucket_policy_document" {
  version = "2012-10-17"
  statement {
    effect = "Allow"
    principals {
      type        = "Service"
      identifiers = ["logdelivery.elasticloadbalancing.amazonaws.com"]
    }
    actions = ["s3:PutObject"]
    resources = [
      "${aws_s3_bucket.alb_log_bucket.arn}/AWSLogs/${data.aws_caller_identity.caller_identity.account_id}/*"
    ]
    condition {
      test     = "ArnLike"
      variable = "aws:SourceArn"
      values   = ["arn:aws:elasticloadbalancing:${local.aws_region}:${data.aws_caller_identity.caller_identity.account_id}:loadbalancer/*"]
    }
  }
}

resource "aws_s3_bucket_policy" "alb_log_bucket_policy" {
  bucket = aws_s3_bucket.alb_log_bucket.id
  policy = data.aws_iam_policy_document.alb_log_bucket_policy_document.json
}

# ALB listener
# We have both http and https listener, http listener should redirect to https listener.
resource "aws_alb_listener" "http_redirect_listener" {
  load_balancer_arn = aws_lb.alb.arn
  protocol          = "HTTP"
  port              = local.http_port

  default_action {
    type = "redirect"
    redirect {
      port        = local.https_port
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }

  tags = {
    Name = "${local.project}-${local.env}-http-redirect-listener"
  }
}

resource "aws_alb_listener" "https_forward_listener" {
  load_balancer_arn = aws_lb.alb.arn
  protocol          = "HTTPS"
  port              = local.https_port
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "arn:aws:acm:ap-northeast-1:957573575820:certificate/4c780941-ca8d-47d2-bc70-219572936118"

  default_action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.target_group.arn
  }

  tags = {
    Name = "${local.project}-${local.env}-https-redirect-listener"
  }
}

resource "aws_alb_target_group" "target_group" {
  name        = "${local.project}-${local.env}-alb-tg"
  vpc_id      = aws_vpc.vpc.id
  target_type = "ip"

  port     = local.http_port
  protocol = "HTTP"

  health_check {
    enabled             = true
    port                = local.http_port
    path                = "/"
    interval            = 30
    healthy_threshold   = 2
    unhealthy_threshold = 2

    matcher = "200-299"
  }

  tags = {
    Name = "${local.project}-${local.env}-alb-target-group"
  }
}
