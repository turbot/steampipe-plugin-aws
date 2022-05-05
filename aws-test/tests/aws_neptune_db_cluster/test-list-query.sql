select db_cluster_identifier, arn, db_cluster_resource_id
from aws_neptune_db_cluster
where arn = '{{ output.resource_aka.value }}';
