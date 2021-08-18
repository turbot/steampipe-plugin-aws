select 
  name, 
  ciphers, 
  ssl_protocols, 
  account_id, 
  partition
from 
  aws.aws_ec2_ssl_policy
where 
  name = '{{ resourceName }}' and region = '{{ output.region_name.value }}';
