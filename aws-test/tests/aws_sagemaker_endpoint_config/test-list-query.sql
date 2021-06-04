select name, arn, partition, region
from aws.aws_sagemaker_endpoint_config
where arn = '{{ output.resource_aka.value }}';