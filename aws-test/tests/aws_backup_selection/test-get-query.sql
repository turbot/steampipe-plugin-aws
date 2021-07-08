select
  backup_plan_id,
  selection_id,
  selection_name,
  arn,
  iam_role_arn,
  resources
from
  aws.aws_backup_selection
where
  backup_plan_id = '{{ output.plan_id.value }}'
  and selection_id = '{{ output.selection_id.value }}';