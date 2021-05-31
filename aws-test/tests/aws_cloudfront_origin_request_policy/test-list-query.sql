select account_id, comment, etag, id, name, title
from aws_cloudfront_origin_request_policy
where akas::text = '["{{ output.resource_aka.value }}"]';