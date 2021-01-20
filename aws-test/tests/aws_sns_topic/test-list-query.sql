select topic_arn, display_name, owner, policy, policy_std, effective_delivery_policy, subscriptions_confirmed, subscriptions_deleted, subscriptions_pending, tags_raw
from aws.aws_sns_topic
where akas::text = '["{{output.resource_aka.value}}"]'
