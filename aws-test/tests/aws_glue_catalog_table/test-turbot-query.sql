select 
  name, 
  region, 
  title
from 
  aws.aws_glue_catalog_table
where 
  name = '{{ resourceName }}';
