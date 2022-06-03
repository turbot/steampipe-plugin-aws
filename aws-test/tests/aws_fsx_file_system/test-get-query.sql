select arn, file_system_id
from aws.aws_fsx_file_system
where file_system_id = '{{ output.resource_id.value }}';