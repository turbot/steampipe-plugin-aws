select akas, tags, region
from aws.aws_cloudfront_distribution
where id = '{{ output.resource_id.value }}';