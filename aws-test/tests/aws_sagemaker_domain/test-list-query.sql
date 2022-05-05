select name, arn, partition, region
from aws.aws_sagemaker_domain
where name = '{{ resourceName }}';