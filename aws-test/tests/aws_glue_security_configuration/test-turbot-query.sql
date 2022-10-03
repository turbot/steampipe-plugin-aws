select 
  akas, 
  name, 
  region, 
  title
from 
  aws.aws_glue_security_configuration
where 
  name = '{{ resourceName }}';
