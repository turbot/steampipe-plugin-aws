select
  title,
  akas
from
  aws_ram_resource_association
where resource_share_name = '{{ resourceName }}';
