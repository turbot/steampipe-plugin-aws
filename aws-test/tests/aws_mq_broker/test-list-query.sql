select arn, engine_type, engine_version, broker_name
from aws.aws_mq_broker
where broker_name = '{{ resourceName }}';
