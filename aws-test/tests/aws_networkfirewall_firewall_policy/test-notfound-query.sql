select 
  arn
from 
  aws.aws_networkfirewall_firewall_policy
where 
  arn = '{{ output.resource_aka.value }}-dummy';
