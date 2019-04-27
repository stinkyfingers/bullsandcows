provider "aws" {
  profile = "jds"
  region = "us-west-1"
}

resource "aws_iam_user" "johnshenk" {
  name = "johnshenk"
}
