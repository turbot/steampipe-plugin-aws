select name, id, control_sources
from aws.aws_auditmanager_control
where name = '{{ output.resource_name.value }}';