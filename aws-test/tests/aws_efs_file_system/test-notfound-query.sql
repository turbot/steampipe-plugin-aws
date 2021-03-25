select title, akas, region, account_id
from aws.aws_efs_file_system
where file_system_id = '{{ output.resource_id.value }}1';