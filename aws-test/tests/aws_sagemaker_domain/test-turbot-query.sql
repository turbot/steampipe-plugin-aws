select title, akas, region, account_id
from aws.aws_sagemaker_domain
where name = '{{ resourceName }}';