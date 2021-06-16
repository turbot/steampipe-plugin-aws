select
  backup_plan_id,
  selection_id,
  arn
from
  aws.aws_backup_selection
where
  akas::text = '["{{ output.resource_aka.value }}"]';