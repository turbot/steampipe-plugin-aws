select name, arn, tags, partition, region, account_id
from aws.aws_sagemaker_notebook_instance
where name = '{{ resourceName }}';
