select db_snapshot_identifier, arn, type, encrypted
from aws.aws_rds_db_snapshot
where db_snapshot_identifier = 'dummy-{{ resourceName }}'
