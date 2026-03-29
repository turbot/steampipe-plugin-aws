select code
from aws.aws_lambda_function
where name = '{{ resourceName }}';
