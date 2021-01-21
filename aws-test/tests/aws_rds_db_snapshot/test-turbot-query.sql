select db_snapshot_identifier, title, tags, akas
from aws.aws_rds_db_snapshot
where db_snapshot_identifier = '{{ resourceName }}'
