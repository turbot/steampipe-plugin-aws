select name, arn, status, function_config
from aws.aws_cloudfront_function
where name = '{{ output.id.value }}';