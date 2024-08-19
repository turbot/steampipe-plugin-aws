select name, arn
from aws_memorydb_cluster
where name = '{{ resourceName }}'
