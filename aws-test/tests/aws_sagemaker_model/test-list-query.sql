select name, arn, partition, region
from aws.aws_sagemaker_model
where arn = '{{ output.resource_aka.value }}';