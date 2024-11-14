resource "aws_dynamodb_table" "basic-dynamodb-table" {
  name           = "GameScores"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Key"

  attribute {
    name = "Key"
    type = "S"
  }

  tags = {
    Name        = var.tag_name
    Environment = var.tag_environment
  }
}
