select name, arn, assume_role_policy, assume_role_policy_std, attached_policy_arns, description, inline_policies, max_session_duration, path, permissions_boundary_arn, permissions_boundary_type, role_last_used_date, role_last_used_region, tags_src, title, tags, akas, partition, account_id
from aws.aws_iam_role
where name = '{{resourceName}}'
