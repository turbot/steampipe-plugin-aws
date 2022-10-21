select
  subnet_group_name,
  description,
  vpc_id,
  region
from
  aws_dax_subnet_group
where subnet_group_name = '{{ resourceName }}-dummy';
