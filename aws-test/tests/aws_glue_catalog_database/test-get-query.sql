select 
  name, 
  catalog_id, 
  create_table_default_permissions, 
  location_uri, 
  description, 
  title, 
  akas
from 
  aws.aws_glue_catalog_database
where 
  name = '{{ resourceName }}';
