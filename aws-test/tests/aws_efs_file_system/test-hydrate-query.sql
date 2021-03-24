select file_system_id, policy, akas, tags, title
from aws.aws_efs_file_system
where file_system_id = '{{ output.resource_id.value }}';