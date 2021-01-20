select akas
from aws.aws_s3_bucket
where name = '{{resourceName}}'
