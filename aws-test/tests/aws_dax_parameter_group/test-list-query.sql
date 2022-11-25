select
  parameter_group_name,
  description,
  region
from
  aws_dax_parameter_group
where
  akas::text = '["{{ output.resource_aka.value }}"]';