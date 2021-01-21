select function_name, arn, version
from aws.aws_lambda_version
where function_name = '{{resourceName}}' and version = '$LATEST'
