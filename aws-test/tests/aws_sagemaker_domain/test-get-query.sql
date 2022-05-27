select name, arn, partition, region, account_id
from aws.aws_sagemaker_domain
where name = '{{ resourceName }}';
