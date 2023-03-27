select
  arn,
  name,
  vpc_id,
  policy_arn,
  region,
  tags
from
  aws.aws_networkfirewall_firewall
where 
  akas = '["{{ output.resource_aka.value }}"]';
