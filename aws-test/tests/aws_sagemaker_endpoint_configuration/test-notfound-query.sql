select title, akas, region, account_id
from aws.aws_sagemaker_endpoint_configuration
where arn = 'dummy-resource';