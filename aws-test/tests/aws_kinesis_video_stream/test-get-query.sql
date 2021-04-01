select stream_name, stream_arn, status, version, kms_key_id, media_type, device_name, data_retention_in_hours, tags, account_id, partition, region
from aws.aws_kinesis_video_stream
where stream_name = '{{ resourceName }}';