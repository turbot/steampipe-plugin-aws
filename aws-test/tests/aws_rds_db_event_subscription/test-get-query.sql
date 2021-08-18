select cust_subscription_id, customer_aws_id, arn, enabled, sns_topic_arn, status
from
  aws.aws_rds_db_event_subscription
where
  cust_subscription_id = '{{ resourceName }}'
