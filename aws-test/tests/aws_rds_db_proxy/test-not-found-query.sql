select 
  db_proxy_name, 
  db_proxy_arn, 
  status
from 
  aws.aws_rds_db_proxy
where 
  db_proxy_name = 'dummy-{{ resourceName }}'
