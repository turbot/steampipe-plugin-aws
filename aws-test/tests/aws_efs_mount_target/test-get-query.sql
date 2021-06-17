select account_id, file_system_id, mount_target_id, title
from aws_efs_mount_target
where mount_target_id = '{{ output.mount_target_id.value }}';