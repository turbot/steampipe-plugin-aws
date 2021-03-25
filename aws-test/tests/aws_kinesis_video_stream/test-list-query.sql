select stream_name, stream_arn, status, version
from aws.aws_kinesis_video_stream
where akas::text = '["{{ output.resource_aka.value }}"]';