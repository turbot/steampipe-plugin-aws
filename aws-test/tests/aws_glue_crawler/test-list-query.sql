select 
  description, 
  name, 
  partition, 
  region, 
  title
from 
  aws.aws_glue_crawler
where 
  akas::text = '["{{ output.resource_aka.value }}"]';
