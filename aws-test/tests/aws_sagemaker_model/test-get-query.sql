select name, arn, partition, region, account_id
from aws.aws_sagemaker_model
where name = '{{ resourceName }}';
