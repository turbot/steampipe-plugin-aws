select tags, akas, title
from aws.aws_sqs_queue
where queue_url = '{{ output.queue_url.value }}'
