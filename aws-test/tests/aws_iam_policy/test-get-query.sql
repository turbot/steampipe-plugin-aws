select name, path, arn, is_attachable, attachment_count, default_version_id, permissions_boundary_usage_count, title
from aws.aws_iam_policy
where arn = '{{ output.resource_aka.value }}'
