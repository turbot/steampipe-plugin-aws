select
  *
from
  aws.aws_cloudfront_response_headers_policy
where
  id = '{{ output.resource_id.value }}';

select
  *
from
  aws.aws_cloudfront_response_headers_policy
where
  id = 'eaab4381-ed33-4a86-88ca-d9558dc6cd63';