select db_cluster_snapshot_arn, db_cluster_identifier, db_cluster_snapshot_identifier, snapshot_type
from aws.aws_neptune_db_cluster_snapshot
where db_cluster_snapshot_arn = '{{ output.resource_aka.value }}'
