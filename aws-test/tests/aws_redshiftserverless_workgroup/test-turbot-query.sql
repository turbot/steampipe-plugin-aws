select
  workgroup_id,
  title,
  tags,
  akas
from
  aws_redshiftserverless_workgroup
where
  workgroup_name = '{{ output.resource_name.value }}';