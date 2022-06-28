connection "aws" {
  plugin = "aws"

  # You may connect to one or more regions. If `regions` is not specified,
  # Steampipe will use a single default region using the same resolution
  # order as the AWS CLI:
  #  1. The `AWS_DEFAULT_REGION` or `AWS_REGION` environment variable
  #  2. The region specified in the active profile (`AWS_PROFILE` or default)
  #regions = ["us-east-1", "us-west-2"]

  # If no credentials are specified, the plugin will use the AWS credentials
  # resolver to get the current credentials in the same manner as the CLI.
  # Alternatively, you may set static credentials with the `access_key`,
  # `secret_key`, and `session_token` arguments, or select a named profile
  # from an AWS credential file with the `profile` argument:
  #profile = "profile2"

  # The maximum number of attempts (including the initial call) Steampipe will
  # make for failing API calls. Can also be set with the AWS_MAX_ATTEMPTS environment variable.
  # Defaults to 9 and must be greater than or equal to 1.
  #max_error_retry_attempts = 9

  # The minimum retry delay in milliseconds after which retries will be performed.
  # This delay is also used as a base value when calculating the exponential backoff retry times.
  # Defaults to 25ms and must be greater than or equal to 1ms.
  #min_error_retry_delay = 25

  # List of additional AWS error codes to ignore for all queries.
  # By default, common not found error codes are ignored and will still be ignored even if this argument is not set.
  #ignore_error_codes = ["AccessDenied", "AccessDeniedException", "NotAuthorized", "UnauthorizedOperation", "UnrecognizedClientException", "AuthorizationError"]

  # Specify the endpoint URL used when making requests to AWS services.
  # If not set, the default AWS generated endpoint will be used.
  # Can also be set with the AWS_ENDPOINT_URL environment variable.
  #endpoint_url = "http://localhost:4566"

  # Set to `true` to force S3 requests to use path-style addressing,
  # i.e., `http://s3.amazonaws.com/BUCKET/KEY`. By default, the S3 client
  # will use virtual hosted bucket addressing when possible (`http://BUCKET.s3.amazonaws.com/KEY`).
  #s3_force_path_style = false
}
