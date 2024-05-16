provider "aws" {

  access_key = var.aws.access_key
  secret_key = var.aws.secret_key
  region     = var.aws.region

  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
  endpoints {
    s3 = "http://s3.localhost.localstack.cloud:4566"
  }
}

resource "aws_s3_bucket" "local-bucket" {
  bucket = "my-bucket"

}
# perform a heuristic check to determine if this has been deployed against localstack
data "aws_caller_identity" "current" {}
output "is_localstack" {
  value = data.aws_caller_identity.current.id == "000000000000"
}
