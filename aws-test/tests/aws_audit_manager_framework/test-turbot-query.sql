select title, akas
from aws.aws_audit_manager_framework
where id = '{{ output.id.value }}';
