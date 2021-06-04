select title, akas, region, account_id
from aws.aws_sagemaker_endpoint_config
where arn = 'dummy-resource';