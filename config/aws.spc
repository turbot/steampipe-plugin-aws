connection "aws" {
  plugin = "aws"

  # `regions` defines the list of regions that Steampipe should target for
  # each query. API calls are made to multiple regions in parallel. The regions
  # list may include wildcards (e.g. *, us-*, us-??st-1).
  # If `regions` is not specified, Steampipe will target the `default_region`
  # only.
  #regions = ["*"] # All regions
  #regions = ["eu-*"] # All EU regions
  #regions = ["us-east-1", "eu-west-2"] # Specific regions

  # Some AWS APIs (e.g. describe EC2 regions, S3 get bucket location) have
  # global results, so can be run against any region. For faster results, you
  # may set your default (closest) region to use for these situations.
  # If not specified, it is resolved in this order:
  #  1. The `AWS_REGION` or `AWS_DEFAULT_REGION` environment variable
  #  2. The region specified in the active profile (`AWS_PROFILE` or default)
  #  3. The main region for the partition as best guessed from the `regions` list.
  #  4. us-east-1, the main region for the most common partition.
  #default_region = "eu-west-2"

  # If no credentials are specified, the plugin will use the AWS credentials
  # resolver to get the current credentials in the same manner as the CLI.
  # Alternatively, you may set static credentials with the `access_key`,
  # `secret_key`, and `session_token` arguments, or select a named profile
  # from an AWS credential file with the `profile` argument:
  #profile = "myprofile"

  # The maximum number of attempts (including the initial call) Steampipe will
  # make for failing API calls. Can also be set with the AWS_MAX_ATTEMPTS environment variable.
  # Defaults to 9 and must be greater than or equal to 1.
  #max_error_retry_attempts = 9

  # The minimum retry delay in milliseconds after which retries will be performed.
  # This delay is also used as a base value when calculating the exponential backoff retry times.
  # Defaults to 25ms and must be greater than or equal to 1ms.
  #min_error_retry_delay = 25

  # List of additional AWS error codes to ignore for all queries.
  # When encountering these errors, the API call will not be retried and empty results will be returned.
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
