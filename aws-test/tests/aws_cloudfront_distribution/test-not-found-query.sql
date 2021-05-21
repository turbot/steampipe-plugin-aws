select akas, tags
from aws.aws_cloudfront_distribution
where id = '{{ output.resource_id.value }}::asd';