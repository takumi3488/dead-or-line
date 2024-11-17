data "aws_iam_policy_document" "event_bridge_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["scheduler.amazonaws.com"]
    }
    effect = "Allow"
  }
}

resource "aws_iam_role" "event_bridge" {
  name               = "role-for-test_lambda-event_bridge"
  assume_role_policy = data.aws_iam_policy_document.event_bridge_assume_role.json
}

resource "aws_iam_role_policy" "event_bridge" {
  name   = "role_policy-for-test_lambda-event_bridge"
  role   = aws_iam_role.event_bridge.name
  policy = data.aws_iam_policy_document.event_bridge.json
}

data "aws_iam_policy_document" "event_bridge" {
  statement {
    effect = "Allow"

    actions = [
      "lambda:InvokeFunction",
    ]

    resources = [
      aws_lambda_function.lambda_function.arn
    ]
  }
}

resource "aws_scheduler_schedule" "test_lambda" {
  name = "test_lambda-event_bridge"

  schedule_expression = "rate(${var.rate_minutes} minutes)"

  flexible_time_window {
    mode = "OFF"
  }

  target {
    arn      = aws_lambda_function.lambda_function.arn
    role_arn = aws_iam_role.event_bridge.arn
  }
}
