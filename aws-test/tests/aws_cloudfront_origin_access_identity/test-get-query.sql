select id, comment, s3_canonical_user_id, etag
from aws.aws_cloudfront_origin_access_identity
where id = '{{ output.resource_id.value }}';