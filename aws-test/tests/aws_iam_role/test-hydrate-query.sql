select name, arn, inline_policies, permissions_boundary_arn, attached_policy_arns, tags_raw
from aws.aws_iam_role
where name = '{{ resourceName }}'
