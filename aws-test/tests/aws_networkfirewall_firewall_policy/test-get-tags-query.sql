select 
  tags_src
from 
  aws.aws_networkfirewall_firewall_policy
where 
  arn = '{{output.resource_aka.value}}';
