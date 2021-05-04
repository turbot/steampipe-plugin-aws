select akas, tags, title
from aws.aws_kinesis_video_stream
where stream_name = '{{ resourceName }}';