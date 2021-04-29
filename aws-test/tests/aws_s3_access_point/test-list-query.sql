select 
  name,
  access_point_arn
from 
  aws.aws_s3_access_point
where 
  access_point_arn = '{{ output.arn.value }}';
