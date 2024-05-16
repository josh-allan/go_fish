variable "aws" {
  description = "AWS Secrets"
  type = object({
    access_key = string
    secret_key = string
    region     = string
  })
}
