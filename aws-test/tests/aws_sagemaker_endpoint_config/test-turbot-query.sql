select title, akas, tags, region, account_id
from aws.aws_sagemaker_endpoint_config
where name = '{{ resourceName }}';