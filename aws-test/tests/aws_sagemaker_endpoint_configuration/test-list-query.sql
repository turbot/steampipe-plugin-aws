select name, arn, partition, region
from aws.aws_sagemaker_endpoint_configuration
where arn = '{{ output.resource_aka.value }}';