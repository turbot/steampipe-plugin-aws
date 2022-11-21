select delivery_stream_name, arn, delivery_stream_status, version_id
from aws_kinesis_firehose_delivery_stream
where akas::text = '["{{ output.resource_aka.value }}"]';