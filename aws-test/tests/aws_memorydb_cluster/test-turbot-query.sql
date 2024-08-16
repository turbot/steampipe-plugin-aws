select title, akas, region, account_id
from aws_memorydb_cluster
where arn = '{{ output.resource_aka.value }}'
