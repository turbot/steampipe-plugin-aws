select db_cluster_snapshot_identifier, title, akas
from aws.aws_docdb_cluster_snapshot
where db_cluster_snapshot_identifier = '{{ resourceName }}'
