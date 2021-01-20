select db_cluster_snapshot_identifier, title, tags, akas
from aws.aws_rds_db_cluster_snapshot
where db_cluster_snapshot_identifier = '{{ resourceName }}'