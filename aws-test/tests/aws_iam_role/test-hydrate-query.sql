select name, arn, inline_policies, permissions_boundary_arn, attached_policy_arns, tags_src
from aws.aws_iam_role
where name = '{{ resourceName }}'
