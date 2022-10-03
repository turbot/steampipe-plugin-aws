select 
  name, 
  connection_type, 
  description,
  region
from 
  aws.aws_glue_connection
where 
  arn = '{{ output.resource_aka.value }}';