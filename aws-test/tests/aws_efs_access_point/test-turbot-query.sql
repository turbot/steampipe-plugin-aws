select akas, tags, title
from aws.aws_efs_access_point
where access_point_id = '{{ output.resource_id.value }}';