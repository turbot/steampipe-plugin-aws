select 
  akas, 
  name, 
  region, 
  title
from 
  aws.aws_glue_catalog_database
where 
  name = '{{ resourceName }}';
