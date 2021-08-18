select cust_subscription_id, arn, status
from aws.aws_rds_db_event_subscription
where cust_subscription_id = 'dummy-{{ resourceName }}';
