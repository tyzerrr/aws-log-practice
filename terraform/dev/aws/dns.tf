# Hosted zone is the DNS record set for this domain.
# The domain's NS records delegate DNS lookups to this Route 53 hosted zone.
data "aws_route53_zone" "main" {
  name         = var.domain_name
  private_zone = false
}

# This A alias record points the apex domain to the ALB.
# Route 53 resolves var.domain_name to the ALB without hard-coding ALB IPs.
resource "aws_route53_record" "alb" {
  zone_id = data.aws_route53_zone.main.zone_id
  name    = var.domain_name
  type    = "A"

  alias {
    # The alias target is the AWS-managed DNS name and zone of the ALB.
    name                   = aws_lb.alb.dns_name
    zone_id                = aws_lb.alb.zone_id
    evaluate_target_health = true
  }
}
