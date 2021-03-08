
select consumer_count, tags_src
from aws.aws_kinesis_stream
where stream_name = '{{ output.resource_name.value }}';