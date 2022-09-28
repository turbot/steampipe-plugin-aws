select
  workgroup_name,
  workgroup_arn,
  workgroup_id,
  region,
  status
from
  aws_redshiftserverless_workgroup
where
  workgroup_name = '{{ output.resource_name.value }}-dummy';