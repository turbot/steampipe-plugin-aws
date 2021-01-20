select db_snapshot_identifier, db_snapshot_attributes
from aws.aws_rds_db_snapshot
where db_snapshot_identifier = '{{ resourceName }}'