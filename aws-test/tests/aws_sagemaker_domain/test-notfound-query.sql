select title, akas, tags, region, account_id
from aws.aws_sagemaker_domain
where name = 'dummy-{{ resourceName }}';