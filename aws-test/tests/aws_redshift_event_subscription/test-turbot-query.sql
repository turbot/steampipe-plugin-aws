select account_id, akas , tags , title
from aws.aws_redshift_event_subscription
where akas::text = '["{{ output.resource_aka.value }}"]';