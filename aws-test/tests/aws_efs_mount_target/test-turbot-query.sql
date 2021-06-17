select account_id, akas, title
from aws_efs_mount_target
where akas::text = '["{{ output.resource_aka.value }}"]';