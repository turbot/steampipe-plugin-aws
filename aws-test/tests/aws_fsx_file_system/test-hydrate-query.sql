select file_system_id, akas, title
from aws.aws_fsx_file_system
where file_system_id = '{{ output.resource_id.value }}';