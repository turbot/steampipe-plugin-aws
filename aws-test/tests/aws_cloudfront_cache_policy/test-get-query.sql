select id, name, comment, e_tag
from aws.aws_cloudfront_cache_policy
where id = '{{ output.resource_id.value }}';