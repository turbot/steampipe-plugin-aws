select id, arn, caller_reference
from aws.aws_cloudfront_origin_access_identity
where id = '{{ output.resource_id.value }}';