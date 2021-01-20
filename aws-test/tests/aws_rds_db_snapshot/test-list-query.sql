select db_snapshot_identifier, arn
from aws.aws_rds_db_snapshot
where arn = '{{ output.resource_aka.value }}'