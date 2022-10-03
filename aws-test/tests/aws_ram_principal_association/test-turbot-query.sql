select
  title,
  akas
from
  aws_ram_principal_association
where resource_share_name = '{{ resourceName }}';
