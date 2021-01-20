select db_cluster_snapshot_identifier, arn, type, db_cluster_identifier
from aws.aws_rds_db_cluster_snapshot
where arn = '{{ output.resource_aka.value }}'