select *
from aws.aws_sns_topic_subscription
where subscription_arn = '{{ output.resource_aka.value }}:dsf'
