select name, id, type, control_sources
from aws.aws_audit_manager_control
where name = '{{ output.resource_name.value }}';