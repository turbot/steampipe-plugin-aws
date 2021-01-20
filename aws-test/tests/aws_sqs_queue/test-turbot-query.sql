select title, akas, tags, region, account_id
from aws.aws_sqs_queue
where queue_url = '{{ output.queue_url.value }}'
