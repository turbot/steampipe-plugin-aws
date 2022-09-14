select
  workgroup_name,
  workgroup_arn,
  workgroup_id,
  region,
  status
from
  aws_redshiftserverless_workgroup
where
  workgroup_arn = '{{ output.resource_aka.value }}';