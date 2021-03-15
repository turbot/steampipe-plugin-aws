select name, arn
from aws.aws_backup_vault
where name = '{{ output.id.value }}'
