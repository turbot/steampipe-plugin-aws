select title, akas
from aws.aws_backup_vault
where name = '{{ output.id.value }}'
