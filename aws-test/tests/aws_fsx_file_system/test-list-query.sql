select file_system_id, title
from aws.aws_fsx_file_system
where akas::text = '["{{ output.resource_aka.value }}"]';