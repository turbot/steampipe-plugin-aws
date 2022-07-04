select title, akas, region, account_id
from aws.aws_sagemaker_app
where name = 'dummy-{{ resourceName }}';