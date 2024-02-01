select arn, broker_name, broker_id
from aws.aws_mq_broker
where broker_id = '{{ output.resource_id.value }}';
