select db_cluster_snapshot_identifier, arn, db_cluster_identifier
from aws.aws_docdb_cluster_snapshot
where db_cluster_snapshot_identifier = 'dummy-{{ resourceName }}'
