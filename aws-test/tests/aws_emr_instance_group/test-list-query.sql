select
  cluster_id,
  instance_group_type,
  instance_type,
  configurations_version,
  last_successfully_applied_configurations_version,
  market,
  requested_instance_count,
  running_instance_count
from
  aws.aws_emr_instance_group
where
  cluster_id = '{{ output.resource_id.value }}'
  and instance_group_type = 'MASTER';