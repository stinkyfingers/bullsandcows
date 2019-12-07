provider "aws" {
  profile = "jds"
  region = "us-west-1"
}

resource "aws_iam_user" "johnshenk" {
  name = "johnshenk"
}

provider "github" {
  token        = "${data.aws_ssm_parameter.github_organization.value}"
  organization = "${data.aws_ssm_parameter.github_token.value}"
}

data "aws_ssm_parameter" "github_organization" {
  name = "/github_organization"
}

data "aws_ssm_parameter" "github_token" {
  name = "/github_token"
}

# cloudfront
locals {
  s3_origin_id = "S3-bullsandcows.john-shenk.com"
}

resource "aws_cloudfront_origin_access_identity" "origin_access_identity" {
  comment = "bullsandccows.john-shenk.com identity"
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  origin {
    domain_name = "${aws_s3_bucket.bullsandcows.bucket_domain_name}"
    origin_id   = "${local.s3_origin_id}"

    s3_origin_config {
      origin_access_identity = "${aws_cloudfront_origin_access_identity.origin_access_identity.cloudfront_access_identity_path}"
    }
  }

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  price_class         = "PriceClass_200"

  aliases = ["bullsandcows.john-shenk.com"]

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "${local.s3_origin_id}"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    compress               = true
    min_ttl                = 0
    default_ttl            = 86400
    max_ttl                = 31536000
  }

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US", "CA"]
    }
  }

  tags = {
    Environment = "production"
  }

  viewer_certificate {
    acm_certificate_arn = "arn:aws:acm:us-east-1:671958020402:certificate/4f5aae32-96c3-4926-ab41-8fb72334470a"
    ssl_support_method = "sni-only"
    minimum_protocol_version = "TLSv1.1_2016"
  }
}

resource "aws_route53_record" "bullsandcows" {
  zone_id = "Z3P68RXJ4VECYX"
  name    = "bullsandcows.john-shenk.com"
  type    = "A"

  alias {
    name                   = "${aws_cloudfront_distribution.s3_distribution.domain_name}"
    zone_id                = "${aws_cloudfront_distribution.s3_distribution.hosted_zone_id}"
    evaluate_target_health = false
  }
}
