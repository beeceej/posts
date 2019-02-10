variable "table_name" {
  description = "name for the dynamo table"
}

variable "hash_key" {
  description = "The hash key of the table"
}

variable "hash_key_type" {
  description = "The type of the hash key"
}

variable "range_key" {
  description = "The range key of the table"
}

variable "range_key_type" {
  description = "The type of the range key"
}

resource "aws_dynamodb_table" "pay-per-request-dynamodb-table" {
  name         = "${var.table_name}"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "${var.hash_key}"
  range_key    = "${var.range_key}"

  attribute {
    name = "${var.hash_key}"
    type = "${var.hash_key_type}"
  }

  attribute {
    name = "${var.range_key}"
    type = "${var.range_key_type}"
  }

  lifecycle {
    prevent_destroy = true
  }
}

output "table_name" {
  value = "${var.table_name}"
}

output "arn" {
  value = "${aws_dynamodb_table.pay-per-request-dynamodb-table.arn}"
}
