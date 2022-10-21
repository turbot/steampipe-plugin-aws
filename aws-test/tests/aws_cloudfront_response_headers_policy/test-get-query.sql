select
  account_id,
  arn,
  etag,
  id,
  name,
  response_headers_policy_config,
  type
from
  aws.aws_cloudfront_response_headers_policy
where
  id = '{{ output.resource_id.value }}';