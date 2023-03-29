select 
  akas,
  db_proxy_name,
  tags,
  title
from 
  aws.aws_rds_db_proxy
where 
  db_proxy_name = '{{ resourceName }}'
