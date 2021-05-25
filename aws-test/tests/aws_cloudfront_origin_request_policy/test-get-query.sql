select account_id, comment, e_tag, id, title
from aws_cloudfront_origin_request_policy
where id = '{{ output.resource_id.value }}';