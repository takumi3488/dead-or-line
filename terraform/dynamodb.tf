resource "aws_dynamodb_table" "table" {
  name           = var.dynamodb_table_name
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Key"

  attribute {
    name = "Key"
    type = "S"
  }

  tags = {
    Name = "${var.tag_prefix}-dynamodb"
  }
}

resource "aws_dynamodb_table_item" "item" {
  table_name = aws_dynamodb_table.table.name
  hash_key   = aws_dynamodb_table.table.hash_key

  item = <<ITEM
  {
    "Key": { "S": "${var.dynamodb_key}" },
    "NotifiedAt": { "N": "0" }
  }
  ITEM
}