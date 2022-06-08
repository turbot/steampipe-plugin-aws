select 
  endpoint_name, 
  arn, 
  partition, 
  region, 
  title
from 
  aws_glue_dev_endpoint
where 
  akas::text = '["{{ output.resource_aka.value }}"]';
