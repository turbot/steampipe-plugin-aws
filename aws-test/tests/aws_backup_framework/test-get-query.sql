select framework_name, arn
from aws.aws_backup_framework
where framework_name = '{{ output.id.value }}';