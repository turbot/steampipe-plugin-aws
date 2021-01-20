select alias_arn, description, function_name, function_version, name, akas, title
from aws.aws_lambda_alias
where name = '{{resourceName}}' and function_name = '{{resourceName}}'