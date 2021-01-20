select topic_arn, endpoint, protocol
from aws.aws_sns_topic_subscription
where subscription_arn = '{{ output.resource_aka.value }}'
