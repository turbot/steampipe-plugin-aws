select name, function_name
from aws.aws_lambda_alias
where name = '{{resourceName}}' and function_name = '{{resourceName}}'