select akas, region, title
from aws.aws_redshift_event_subscription
where akas::text = 'abpoxsc';