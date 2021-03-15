select stream_name, stream_arn, stream_status
from aws.aws_kinesis_stream
where akas::text = '["{{ output.resource_aka.value }}"]';