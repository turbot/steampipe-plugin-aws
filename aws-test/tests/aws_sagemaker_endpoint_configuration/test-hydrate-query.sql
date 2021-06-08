select name, tags
from aws.aws_sagemaker_endpoint_configuration
where name = '{{ resourceName }}';
