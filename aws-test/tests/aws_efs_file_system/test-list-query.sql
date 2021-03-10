select akas, file_system_id, encrypted, performance_mode, title
from aws.aws_efs_file_system
where akas::text = '["{{ output.resource_aka.value }}"]';