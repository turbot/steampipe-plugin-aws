select akas, title, region
from aws.aws_cloudfront_cache_policy
where id = '{{ output.resource_id.value }}';