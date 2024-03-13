select
    akas,
    title,
    region
from
    aws.aws_cloudfront_response_headers_policy
where
    id = '{{ output.resource_id.value }}';