select group_id, name, inline_policies, attached_policy_arns, users
from aws.aws_iam_group
where name = '{{resourceName}}'
