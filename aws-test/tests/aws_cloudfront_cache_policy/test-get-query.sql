select id, name, comment, etag
from aws.aws_cloudfront_cache_policy
where id = '{{ output.resource_id.value }}';