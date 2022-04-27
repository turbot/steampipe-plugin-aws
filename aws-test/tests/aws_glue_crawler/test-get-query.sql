select 
  name, 
  akas, 
  description,
  title
from 
  aws.aws_glue_crawler
where 
  name = '{{ resourceName }}';
