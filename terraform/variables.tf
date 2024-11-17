# タグの設定
variable "tag_prefix" {
  type = string
}

# Lambdaの設定
variable "target_url" {
  type = string
}
variable "notice_interval" {
  type = number
}
variable "dynamodb_table_name" {
  type = string
}
variable "dynamodb_key" {
  type = string
}
variable "line_token" {
  type = string
}
variable "line_base_message" {
  type = string
}
variable "line_to" {
  type = string
}
variable "rate_minutes" {
  type = number
}