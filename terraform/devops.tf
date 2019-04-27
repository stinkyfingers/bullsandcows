resource "aws_s3_bucket" "bullsandcows" {
  bucket = "bullsandcows.john-shenk.com"
  acl    = ""
  force_destroy = false
  website {
    index_document = "index.html"
  }
}

resource "aws_s3_bucket_policy" "bullsandcows" {
  bucket = "${aws_s3_bucket.bullsandcows.id}"
  policy = <<EOF
{
  "Version":"2012-10-17",
  "Statement":[
    {"Sid":"PublicReadForGetBucketObjects",
      "Effect":"Allow","Principal":"*",
      "Action":"s3:GetObject",
      "Resource":"arn:aws:s3:::bullsandcows.john-shenk.com/*"
    }]
}
EOF
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
        "Service": "codepipeline.amazonaws.com"
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
        "codepipeline:*"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_codepipeline" "bullsandcows" {
  name     = "bullsandcows"
  role_arn = "${aws_iam_role.bullsandcows.arn}"

  artifact_store {
    location = "${aws_s3_bucket.bullsandcows.bucket}"
    type     = "S3"
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "ThirdParty"
      provider         = "GitHub"
      version          = "1"
      output_artifacts = ["source_output"]

      configuration = {
        Owner  = "stinkyfingers"
        Repo   = "bullsandcows"
        Branch = "master"
        OAuthToken = ""
      }
    }
  }

  stage {
    name = "Build"

    action {
      name            = "Build"
      category        = "Build"
      owner           = "AWS"
      provider        = "CodeBuild"
      input_artifacts = ["source_output"]
      output_artifacts = ["build_output"]
      version         = "1"

      configuration = {
        ProjectName = "bullsandcows"
      }
    }
  }
  #
  # stage {
  #   name = "Deploy"
  #
  #   action {
  #     name             = "Deploy"
  #     category         = "Deploy"
  #     owner            = "AWS"
  #     provider         = "CloudFormation"
  #     input_artifacts  = ["build_output"]
  #     version          = "1"
  #
  #     configuration {
  #       ActionMode     = "REPLACE_ON_FAILURE"
  #       Capabilities   = "CAPABILITY_AUTO_EXPAND,CAPABILITY_IAM"
  #       OutputFileName = "CreateStackOutput.json"
  #       StackName      = "MyStack"
  #       TemplatePath   = "build_output::sam-templated.yaml"
  #     }
  #   }
  # }
}
