select stream_name, stream_arn, stream_status, encryption_type, key_id, retention_period_hours, account_id, partition, region
from aws.aws_kinesis_stream
where stream_name = '{{ resourceName }}';