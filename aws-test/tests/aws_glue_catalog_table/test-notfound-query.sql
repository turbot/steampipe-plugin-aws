select 
  title, 
  akas, 
  region, 
  account_id
from 
  aws.aws_glue_catalog_table
where 
  name = '{{ resourceName }}::xzq';
