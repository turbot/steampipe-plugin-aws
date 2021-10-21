select name, file_system_id
from aws.aws_fsx_file_system
where akas::text = '["{{ output.resource_aka.value }}"]';