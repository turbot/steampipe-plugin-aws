select name, arn
from aws.aws_backup_vault
where akas::text = '["{{ output.resource_aka.value }}"]';
