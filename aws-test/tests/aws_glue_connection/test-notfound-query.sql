select 
  name, 
  connection_type, 
  description
from 
  aws.aws_glue_connection
where 
  name = 'dummy-{{ resourceName }}';