select name, tags
from aws.aws_sagemaker_endpoint_config
where name = '{{ resourceName }}';
