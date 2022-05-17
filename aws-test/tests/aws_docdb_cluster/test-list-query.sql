select db_cluster_identifier, arn
from aws.aws_docdb_cluster
where arn = '{{ output.resource_aka.value }}'
