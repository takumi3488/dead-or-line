data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_cloudwatch_log_group" "lambda" {
  name              = "/aws/lambda/${aws_lambda_function.lambda_function.function_name}"
  retention_in_days = 3
  skip_destroy      = false
}

data "aws_iam_policy_document" "lambda_dynamodb" {
  statement {
    effect = "Allow"
    actions = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
      "dynamodb:UpdateItem",
    ]
    resources = [
      aws_dynamodb_table.table.arn
    ]
  }
}

resource "aws_iam_policy" "lambda_dynamodb" {
  name   = "dead-or-line-lambda-dynamodb-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.lambda_dynamodb.json
}

resource "aws_iam_role_policy_attachment" "lambda_role_attachment" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_dynamodb.arn
}

resource "aws_lambda_function" "lambda_function" {
  function_name = "dead-or-line-function"
  role          = aws_iam_role.iam_for_lambda.arn
  architectures = ["arm64"]
  package_type  = "Image"
  image_uri     = "${aws_ecr_repository.repo.repository_url}@${docker_registry_image.lambda.sha256_digest}"
  memory_size   = 128
  timeout       = 30
  environment {
    variables = {
      TARGET_URL          = var.target_url
      NOTICE_INTERVAL     = var.notice_interval
      DYNAMODB_TABLE_NAME = var.dynamodb_table_name
      DYNAMODB_KEY        = var.dynamodb_key
      LINE_TOKEN          = var.line_token
      LINE_BASE_MESSAGE   = var.line_base_message
      LINE_TO             = var.line_to
    }
  }
  tags = {
    Name = "${var.tag_prefix}-lambda"
  }
}
