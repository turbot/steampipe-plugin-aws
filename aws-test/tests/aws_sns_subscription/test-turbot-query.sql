select title, akas, account_id, region, partition
from aws.aws_sns_subscription
where subscription_arn = '{{ output.resource_aka.value }}'
