select alias_arn, description, function_name, function_version, name, akas, title
from aws.aws_lambda_alias
where akas::text = '["{{ output.resource_aka.value }}"]'