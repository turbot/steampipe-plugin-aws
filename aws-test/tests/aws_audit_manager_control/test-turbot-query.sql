select title, tags, akas
from aws.aws_audit_manager_control
where id = '{{ output.control_id.value }}';