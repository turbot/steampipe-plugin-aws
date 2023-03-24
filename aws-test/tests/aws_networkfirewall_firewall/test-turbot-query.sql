select 
  title, 
  akas, 
  region, 
  account_id
from
  aws.aws_networkfirewall_firewall
where 
  arn = '{{ output.resource_aka.value }}';
