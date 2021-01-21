select topic_arn, subscription_arn, owner, protocol, endpoint
from aws.aws_sns_topic_subscription
where akas::text = '["{{output.resource_aka.value}}"]'
