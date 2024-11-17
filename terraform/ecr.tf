data "aws_ecr_authorization_token" "token" {
}

resource "aws_ecr_repository" "repo" {
  name         = "dead-or-line-repo"
  force_delete = true
}
