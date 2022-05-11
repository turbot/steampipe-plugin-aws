select 
  name, 
  catalog_id, 
  description
from 
  aws.aws_glue_catalog_table
where 
  name = '{{ resourceName }}' 
  and database_name = '{{ output.database_name.value }}';
