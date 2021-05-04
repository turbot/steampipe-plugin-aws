select name, arn, partition, region
from aws.aws_sagemaker_notebook_instance
where name = '{{ resourceName }}';