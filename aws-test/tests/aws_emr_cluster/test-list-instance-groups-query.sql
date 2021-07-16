select
  distinct name,
  id,
  ig ->> 'InstanceType' as instance_type
from
  aws_emr_cluster,
  jsonb_array_elements(instance_groups) as ig
where
  id = '{{ output.resource_id.value }}'
  and ig ->> 'InstanceGroupType' = 'MASTER';