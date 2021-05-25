select comment, id, e_tag, name
from aws.aws_cloudfront_cache_policy
where name = '{{ resourceName }}';