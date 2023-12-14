select account_id, akas, region, tags, title
from aws.aws_mq_broker
where broker_name = 'dummy';
