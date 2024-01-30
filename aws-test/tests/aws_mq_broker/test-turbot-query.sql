select akas, broker_name, region, tags, title
from aws.aws_mq_broker
where broker_name = '{{ resourceName }}';
