select name, id, type, control_sources
from aws.aws_auditmanager_control
where name = '{{ output.resource_name.value }}';