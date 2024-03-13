select title, akas
from aws.aws_auditmanager_framework
where arn = '{{ output.arn.value }}';
