select name, access_point_id, access_point_arn, life_cycle_state, file_system_id, owner_id, posix_user, root_directory, tags_src
from aws.aws_efs_access_point
where access_point_id = '{{ output.resource_id.value }}';