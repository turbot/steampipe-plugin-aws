select account_id, akas, title
from aws_cloudfront_origin_request_policy
where akas::text = '["{{ output.resource_aka.value }}"]';