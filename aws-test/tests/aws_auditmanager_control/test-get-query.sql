select name, id, type, control_sources
from aws.aws_auditmanager_control
where id = '{{ output.control_id.value }}' and region = '{{ output.aws_region.value }}';