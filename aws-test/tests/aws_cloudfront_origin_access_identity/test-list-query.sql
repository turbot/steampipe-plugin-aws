select id, arn, comment, etag
from aws.aws_cloudfront_origin_access_identity
where akas::text = '["{{ output.resource_aka.value }}"]'