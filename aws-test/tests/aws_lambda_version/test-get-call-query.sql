select function_name, arn, version, runtime, akas, title
from aws.aws_lambda_version
where function_name = '{{resourceName}}' and version = '$LATEST'