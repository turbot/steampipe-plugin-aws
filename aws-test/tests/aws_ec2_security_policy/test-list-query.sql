select 
  name, 
  title
from 
  aws.aws_ec2_security_policy
where 
  akas::text = '["arn:{{ output.aws_partition.value }}:elbv2:{{ output.region_name.value }}:{{ output.account_id.value }}:ssl-policy/{{ resourceName }}"]';
