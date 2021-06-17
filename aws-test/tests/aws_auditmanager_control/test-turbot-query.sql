select title, tags, akas
from aws.aws_auditmanager_control
where id = '{{ output.control_id.value }}';