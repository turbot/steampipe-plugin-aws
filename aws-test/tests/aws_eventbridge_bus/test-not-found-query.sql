select account_id, akas, region, tags, title
from aws.aws_eventbridge_bus
where name = 'dummy';
 