select 
  name,
  access_point_policy_is_public
from 
  aws.aws_s3_access_point
where 
  name = '{{ resourceName }}';
