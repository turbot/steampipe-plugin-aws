select delivery_stream_name, arn, delivery_stream_status, delivery_stream_type, version_id, last_update_timestamp, delivery_stream_encryption_configuration, source, failure_description, has_more_destinations, tags_src, region, partition, account_id
from aws_kinesis_firehose_delivery_stream
where delivery_stream_name = '{{ resourceName }}';