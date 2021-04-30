select 
  name, 
  ciphers, 
  ssl_protocols, 
  account_id, 
  partition
from 
  aws.aws_ec2_security_policy
where 
  name = '{{ resourceName }}';
