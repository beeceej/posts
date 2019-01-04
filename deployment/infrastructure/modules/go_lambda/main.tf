variable "handler" {
  description = "handler for lambda function"
}

variable "function_name" {
  description = "name for lambda function"
}

variable "file_name" {
  description = "name for lambda function"
}

variable "environment_vars" {
  description = "environment variables for lambda"
  type        = "map"

  default = {
    state_machine = "blog-post-pipeline"
  }
}

variable "memory_size" {
  default     = 128
  description = "amount of memory capacity for the lambda function"
}

variable "timeout" {
  default     = 60
  description = "amount of time for the lambda to run before being killed"
}

data aws_iam_policy_document "lambda_assume_role" {
  statement {
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "lambda_iam_role" {
  name = "${var.function_name}"

  assume_role_policy = "${data.aws_iam_policy_document.lambda_assume_role.json}"
}

resource "aws_iam_role_policy_attachment" "role_attachment" {
  role       = "${aws_iam_role.lambda_iam_role.name}"
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "state_machine_lambda" {
  function_name    = "${var.function_name}"
  role             = "${aws_iam_role.lambda_iam_role.arn}"
  handler          = "${var.handler}"
  filename         = "${var.file_name}"
  source_code_hash = "${base64sha256(file(var.file_name))}"
  runtime          = "go1.x"
  memory_size      = "${var.memory_size}"
  timeout          = "${var.timeout}"

  environment {
    variables = "${var.environment_vars}"
  }
}

output "lambda_arn" {
  value = "${aws_lambda_function.state_machine_lambda.arn}"
}

output "lambda_role_name" {
  value = "${aws_iam_role.lambda_iam_role.name}"
}
