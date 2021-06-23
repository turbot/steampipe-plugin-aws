select
  title,
  akas
from
  aws.aws_backup_selection
where
  backup_plan_id = '{{ output.plan_id.value }}'
  and selection_id = '{{ output.selection_id.value }}';