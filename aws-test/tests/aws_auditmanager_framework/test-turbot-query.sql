select title, akas
from aws.aws_auditmanager_framework
where id = '{{ output.id.value }}';
