select
  parameter_name,
  description,
  region
from
  aws_dax_parameter
where
  parameter_group_name = '{{ resourceName }}-dummy';
