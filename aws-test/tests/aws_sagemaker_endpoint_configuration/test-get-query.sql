select name, arn, partition, region, account_id
from aws.aws_sagemaker_endpoint_configuration
where name = '{{ resourceName }}';
