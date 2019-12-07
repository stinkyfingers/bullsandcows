resource "aws_s3_bucket" "bullsandcows" {
  bucket = "bullsandcows.john-shenk.com"
  acl    = "private"
  force_destroy = false
  website {
    index_document = "index.html"
  }
}

resource "aws_s3_bucket_policy" "bullsandcows" {
  bucket = "${aws_s3_bucket.bullsandcows.id}"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Cloudfront Read",
            "Effect": "Allow",
            "Principal": {
                "AWS": "${aws_cloudfront_origin_access_identity.origin_access_identity.iam_arn}"
            },
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::bullsandcows.john-shenk.com/*"
        }
    ]
}
EOF
}

resource "aws_s3_bucket_public_access_block" "example" {
  bucket = "${aws_s3_bucket.bullsandcows.id}"
  block_public_acls   = true
  block_public_policy = true
}

resource "aws_iam_role" "bullsandcows" {
  name = "bullsandcows_pipeline_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "bullsandcows" {
  name = "bullsandcows_pipeline_policy"
  role = "${aws_iam_role.bullsandcows.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect":"Allow",
      "Action": [
        "s3:Get*",
        "s3:List*",
        "s3:PutObject"
      ],
      "Resource": [
        "${aws_s3_bucket.bullsandcows.arn}",
        "${aws_s3_bucket.bullsandcows.arn}/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "codebuild:BatchGetBuilds",
        "codebuild:StartBuild"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "codebuild:*"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
    }
  ]
}
EOF
}


resource "aws_codebuild_project" "bullsandcows" {
  name          = "bullsandcows"
  description   = "bullsandcows"
  build_timeout = "5"
  service_role  = "${aws_iam_role.bullsandcows.arn}"

  artifacts {
    type = "NO_ARTIFACTS"
  }

  cache {
    type     = "S3"
    location = "${aws_s3_bucket.bullsandcows.bucket}"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "aws/codebuild/standard:2.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"

    # environment variables go here
  }

  source {
    type            = "GITHUB"
    location        = "https://github.com/stinkyfingers/bullsandcows.git"
    git_clone_depth = 1
    buildspec       = "buildspec.yml"
  }

  tags = {
    "Environment" = "Prod"
  }
}

resource "aws_codebuild_webhook" "bullsandcows" {
  project_name = "${aws_codebuild_project.bullsandcows.name}"
}


# resource "aws_iam_role" "codebuild-bullsandcows-service-role" {
#   name = "codebuild-bullsandcows-service-role"
#   assume_role_policy = <<EOF
# {
#   "Version":"2012-10-17",
#   "Statement":[
#     {
#       "Effect":"Allow",
#       "Principal":{
#         "Service":"codebuild.amazonaws.com"
#       },
#       "Action":"sts:AssumeRole"
#     }
#   ]
# }
# EOF
# }
#
# resource "aws_iam_role_policy" "bullsandcows-service-policy" {
#   name = "bullsandcows_service_policy"
#   # role = "${aws_iam_role.codebuild-bullsandcows-service-role.id}"
#   role = "codebuild-bullsandcows-service-role"
#
#   policy = <<EOF
# {
#   "Version": "2012-10-17",
#   "Statement": [
#     {
#       "Effect":"Allow",
#       "Action": [
#         "s3:Get*",
#         "s3:List*",
#         "s3:PutObject"
#       ],
#       "Resource": [
#         "${aws_s3_bucket.bullsandcows.arn}",
#         "${aws_s3_bucket.bullsandcows.arn}/*"
#       ]
#     }
#   ]
# }
# EOF
# }
#
# resource "aws_codepipeline" "bullsandcows" {
#   name     = "bullsandcows"
#   role_arn = "${aws_iam_role.bullsandcows.arn}"
#
#   artifact_store {
#     location = "${aws_s3_bucket.bullsandcows.bucket}"
#     type     = "S3"
#   }
#
#   stage {
#     name = "Source"
#
#     action {
#       name             = "Source"
#       category         = "Source"
#       owner            = "ThirdParty"
#       provider         = "GitHub"
#       version          = "1"
#       output_artifacts = ["source_output"]
#
#       configuration = {
#         Owner  = "stinkyfingers"
#         Repo   = "bullsandcows"
#         Branch = "master"
#         OAuthToken = "${var.oauthtoken}"
#       }
#     }
#   }
#
#   stage {
#     name = "Build"
#
#     action {
#       name            = "Build"
#       category        = "Build"
#       owner           = "AWS"
#       provider        = "CodeBuild"
#       input_artifacts = ["source_output"]
#       output_artifacts = ["build_output"]
#       version         = "1"
#
#       configuration = {
#         ProjectName = "bullsandcows"
#       }
#     }
#   }
#   #
#   # stage {
#   #   name = "Deploy"
#   #
#   #   action {
#   #     name             = "Deploy"
#   #     category         = "Deploy"
#   #     owner            = "AWS"
#   #     provider         = "CloudFormation"
#   #     input_artifacts  = ["build_output"]
#   #     version          = "1"
#   #
#   #     configuration {
#   #       ActionMode     = "REPLACE_ON_FAILURE"
#   #       Capabilities   = "CAPABILITY_AUTO_EXPAND,CAPABILITY_IAM"
#   #       OutputFileName = "CreateStackOutput.json"
#   #       StackName      = "MyStack"
#   #       TemplatePath   = "build_output::sam-templated.yaml"
#   #     }
#   #   }
#   # }
# }
#
# locals {
#   webhook_secret = "super-secret"
# }
#
# resource "aws_codepipeline_webhook" "bullsandcows" {
#   name            = "bullsandcows-webhook"
#   authentication  = "GITHUB_HMAC"
#   target_action   = "Source"
#   target_pipeline = "${aws_codepipeline.bullsandcows.name}"
#
#   authentication_configuration {
#     secret_token = "${local.webhook_secret}"
#   }
#
#   filter {
#     json_path    = "$.ref"
#     match_equals = "refs/heads/master"
#   }
# }
#
# # Wire the CodePipeline webhook into a GitHub repository.
# resource "github_repository_webhook" "bullsandcows" {
#   repository = "bullsandcows"
#
#   name = "web"
#
#   configuration {
#     url          = "${aws_codepipeline_webhook.bullsandcows.url}"
#     content_type = "form"
#     insecure_ssl = true
#     secret       = "${local.webhook_secret}"
#   }
#
#   events = ["push"]
# }
