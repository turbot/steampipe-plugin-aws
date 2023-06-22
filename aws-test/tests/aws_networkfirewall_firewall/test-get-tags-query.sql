select 
  tags_src
from 
  aws.aws_networkfirewall_firewall
where 
  arn = '{{output.resource_aka.value}}';
