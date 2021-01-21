select function_name, arn, version, runtime, akas, title
from aws.aws_lambda_version
where akas = '["{{ output.resource_aka.value }}:$LATEST"]'
