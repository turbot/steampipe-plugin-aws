select name, arn, inline_policies, attached_policy_arns, permissions_boundary_arn, permissions_boundary_type, tags_raw
from aws.aws_iam_user
where name = '{{ resourceName }}'
