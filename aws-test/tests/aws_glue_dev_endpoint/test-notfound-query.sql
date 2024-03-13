select 
  endpoint_name, 
  arn, 
  region, 
  account_id
from 
  aws_glue_dev_endpoint
where 
  endpoint_name = '{{ resourceName }}::xzq';
