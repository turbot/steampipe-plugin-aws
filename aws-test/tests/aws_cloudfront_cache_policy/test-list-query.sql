select comment, id, etag, name
from aws.aws_cloudfront_cache_policy
where name = '{{ resourceName }}';