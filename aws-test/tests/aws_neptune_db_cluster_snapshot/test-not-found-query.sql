select db_cluster_snapshot_identifier, db_cluster_snapshot_arn, snapshot_type, db_cluster_identifier
from aws.aws_neptune_db_cluster_snapshot
where db_cluster_snapshot_identifier = 'dummy-{{ resourceName }}'
