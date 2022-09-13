select title, akas
from aws.aws_cloudfront_function
where name = '{{ output.id.value }}';