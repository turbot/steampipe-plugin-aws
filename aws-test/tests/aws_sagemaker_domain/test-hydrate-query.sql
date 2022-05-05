select name, akas
from aws.aws_sagemaker_domain
where name = '{{ resourceName }}';
