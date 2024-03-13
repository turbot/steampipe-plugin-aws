select
  resource_share_name,
  resource_share_arn,
  associated_entity
from
  aws_ram_principal_association
where resource_share_name = '{{ resourceName }}';
