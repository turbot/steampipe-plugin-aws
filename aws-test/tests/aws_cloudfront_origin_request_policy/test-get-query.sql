select account_id, comment, etag, id, title
from aws_cloudfront_origin_request_policy
where id = '{{ output.resource_id.value }}';