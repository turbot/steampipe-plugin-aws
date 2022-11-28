select
  parameter_group_name,
  description,
  region
from
  aws_dax_parameter_group
where
  parameter_group_name = '{{ resourceName }}';