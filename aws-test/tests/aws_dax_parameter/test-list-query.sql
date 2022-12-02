select
  parameter_group_name,
  parameter_name,
  parameter_value,
  region
from
  aws_dax_parameter
where
  parameter_group_name = '{{ output.parameter_group_name.value }}' and parameter_name = 'query-ttl-millis';