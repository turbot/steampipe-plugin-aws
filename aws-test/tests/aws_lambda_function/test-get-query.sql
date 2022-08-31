select name, arn, description, version, role
from aws.aws_lambda_function
where name = '{{ resourceName }}';
