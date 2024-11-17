terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.0.2"
    }
  }
}

provider "docker" {
  registry_auth {
    address  = data.aws_ecr_authorization_token.token.proxy_endpoint
    username = data.aws_ecr_authorization_token.token.user_name
    password = data.aws_ecr_authorization_token.token.password
  }
}

resource "docker_image" "lambda" {
  name         = "${aws_ecr_repository.repo.repository_url}:latest"
  platform     = "linux/${var.architecture}"
  keep_locally = true
  build {
    context = "${path.module}/../"
  }
  triggers = {
    dir_sha1 = sha1(join("", [for f in fileset(path.module, "../*") : filesha1(f)]))
  }
}

resource "docker_registry_image" "lambda" {
  name          = docker_image.lambda.name
  keep_remotely = true
  triggers = {
    dir_sha1 = sha1(join("", [for f in fileset(path.module, "../*") : filesha1(f)]))
  }
}
