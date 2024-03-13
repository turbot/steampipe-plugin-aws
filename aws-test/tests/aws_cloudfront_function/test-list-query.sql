select name, arn, status, function_config
from aws.aws_cloudfront_function
where akas::text = '["{{ output.resource_aka.value }}"]';