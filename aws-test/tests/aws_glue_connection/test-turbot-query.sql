select 
  akas, 
  name, 
  region, 
  title
from 
  aws.aws_glue_connection
where 
  name = '{{ resourceName }}';
