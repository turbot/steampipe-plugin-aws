select 
  description, 
  name, 
  region
from 
  aws.aws_glue_catalog_table
where 
  name = '{{ resourceName }}';
