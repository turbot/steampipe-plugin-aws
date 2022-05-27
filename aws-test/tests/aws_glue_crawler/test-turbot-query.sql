select 
  akas, 
  name, 
  region, 
  title
from 
  aws.aws_glue_crawler
where 
  name = '{{ resourceName }}';
