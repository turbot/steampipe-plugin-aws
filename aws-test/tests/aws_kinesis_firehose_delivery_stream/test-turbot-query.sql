select akas, tags, title
from aws.aws_kinesis_firehose_delivery_stream
where delivery_stream_name = '{{ resourceName }}';