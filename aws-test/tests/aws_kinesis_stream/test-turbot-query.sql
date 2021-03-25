select akas, tags, title
from aws.aws_kinesis_stream
where stream_name = '{{ resourceName }}';