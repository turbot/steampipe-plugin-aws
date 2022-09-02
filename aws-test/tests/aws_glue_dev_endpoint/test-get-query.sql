select 
  endpoint_name, 
  arn, 
  status,
  title
from 
  aws_glue_dev_endpoint
where 
  endpoint_name = '{{ resourceName }}';
