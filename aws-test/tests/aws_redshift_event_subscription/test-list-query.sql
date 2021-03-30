select account_id, cust_subscription_id, partition, region, title
from aws.aws_redshift_event_subscription
where akas::text = '["{{ output.resource_aka.value }}"]';