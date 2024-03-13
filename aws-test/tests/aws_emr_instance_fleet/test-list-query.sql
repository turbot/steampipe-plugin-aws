select
  cluster_id,
  instance_fleet_type,
  provisioned_on_demand_capacity
from
  aws.aws_emr_instance_fleet
where
  cluster_id = '{{ output.resource_id.value }}';

