select 
  title, 
  akas, 
  region, 
  account_id
from 
  aws.aws_glue_crawler
where 
  name = '{{ resourceName }}::xzq';
