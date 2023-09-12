select db_cluster_snapshot_identifier, db_cluster_snapshot_attributes
from aws.aws_neptune_db_cluster_snapshot
where db_cluster_snapshot_identifier = '{{ resourceName }}'
