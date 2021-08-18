select cust_subscription_id, arn, enabled
from aws.aws_rds_db_event_subscription
where arn = '{{ output.resource_aka.value }}'
