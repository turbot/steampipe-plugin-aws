select arn, name
from aws_memorydb_cluster
where akas::text = '["{{ output.resource_aka.value }}"]'
