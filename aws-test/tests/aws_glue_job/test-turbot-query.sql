select 
  akas, 
  name, 
  region, 
  title
from 
  aws_glue_job
where 
  name = '{{ resourceName }}';
