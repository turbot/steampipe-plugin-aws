select file_system_id, encrypted, performance_mode, tags_src
from aws.aws_efs_file_system
where file_system_id = '{{ output.resource_id.value }}';