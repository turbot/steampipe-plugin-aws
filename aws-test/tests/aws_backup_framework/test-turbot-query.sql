select akas
from aws.aws_backup_framework
where framework_name = '{{ output.id.value }}';