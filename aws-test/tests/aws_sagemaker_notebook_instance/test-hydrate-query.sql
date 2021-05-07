select name, akas, tags
from aws.aws_sagemaker_notebook_instance
where name = '{{ resourceName }}';
