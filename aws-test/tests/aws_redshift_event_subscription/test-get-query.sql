select account_id, cust_subscription_id, partition, region, tags_src, title
from aws.aws_redshift_event_subscription
where cust_subscription_id = '{{ resourceName}}';