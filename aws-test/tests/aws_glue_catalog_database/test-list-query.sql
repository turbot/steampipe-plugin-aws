select 
  description, 
  name, 
  partition, 
  region, 
  title
from 
  aws.aws_glue_catalog_database
where 
  akas::text = '["{{ output.resource_aka.value }}"]';
