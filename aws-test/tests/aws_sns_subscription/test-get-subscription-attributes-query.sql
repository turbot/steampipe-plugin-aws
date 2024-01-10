select confirmation_was_authenticated, delivery_policy, effective_delivery_policy, redrive_policy, filter_policy, pending_confirmation, raw_message_delivery
from aws.aws_sns_subscription
where subscription_arn = '{{ output.resource_aka.value }}'
