select db_cluster_identifier, arn, resource_id
from aws.aws_rds_db_cluster
where arn = '{{ output.resource_aka.value }}'