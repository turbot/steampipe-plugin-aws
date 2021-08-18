select cust_subscription_id, title, akas
from aws.aws_rds_db_event_subscription
where cust_subscription_id = '{{ resourceName }}'
