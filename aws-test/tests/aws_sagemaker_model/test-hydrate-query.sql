select name, akas
from aws.aws_sagemaker_model
where name = '{{ resourceName }}';
