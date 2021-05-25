select account_id, comment, e_tag, id, name, title
from aws_cloudfront_origin_request_policy
where akas::text = '["{{ output.resource_aka.value }}"]';