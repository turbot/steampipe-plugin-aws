select account_id, file_system_id, mount_target_id, partition, title
from aws_efs_mount_target
where file_system_id = '{{ output.file_system_id.value }}';