select continuous_backups_status, point_in_time_recovery_description
from aws.aws_dynamodb_table
where name = '{{ resourceName }}'
