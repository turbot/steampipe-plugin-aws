select queue_arn, queue_url, message_retention_seconds, max_message_size, delay_seconds, receive_wait_time_seconds, partition, policy, policy_std, tags
from aws.aws_sqs_queue
where akas::text = '["{{output.resource_aka.value}}"]'
