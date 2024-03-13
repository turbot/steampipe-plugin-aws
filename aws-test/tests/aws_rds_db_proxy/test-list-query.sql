select 
  db_proxy_name,
  db_proxy_arn
from 
  aws.aws_rds_db_proxy
where 
  db_proxy_arn = '{{ output.resource_aka.value }}'


