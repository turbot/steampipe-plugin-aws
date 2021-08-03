connection "aws" {
  plugin    = "aws"

  # You may connect to one or more regions. If `regions` is not specified,
  # Steampipe will use a single default region using the same resolution
  # order as the aws cli:
  #  1. The `AWS_DEFAULT_REGION` or `AWS_REGION` environment variable
  #  2. The region specified in the active profile (`AWS_PROFILE` or default)
  #regions     = ["us-east-1", "us-west-2"]

  # If no credentials are specified, the plugin will use the AWS credentials
  # resolver to get the current credentials in the same manner as the CLI
  #  Alternatively, you may set static credentials with the `access_key`, 
  # `secret_key`, and `session_token` arguments, or select a named profile
  # from an AWS credential file with the `profile` argument:
  #profile     = "profile2"
}
