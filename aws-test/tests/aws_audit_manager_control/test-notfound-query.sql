select name, id, type
from aws.aws_audit_manager_control
where name = '{{ output.resource_name.value }}ddd';