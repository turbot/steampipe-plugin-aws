select framework_name, arn
from aws.aws_backup_framework
where akas::text = '["{{ output.resource_aka.value }}"]';