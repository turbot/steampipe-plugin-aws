select 
  akas, 
  endpoint_name, 
  region, 
  title
from 
  aws_glue_dev_endpoint
where 
  endpoint_name = '{{ resourceName }}';
