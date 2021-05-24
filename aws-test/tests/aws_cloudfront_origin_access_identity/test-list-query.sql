select id, arn, comment, e_tag
from aws.aws_cloudfront_origin_access_identity
where akas::text = '["{{ output.resource_aka.value }}"]'