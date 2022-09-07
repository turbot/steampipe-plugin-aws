select 
  application_id,
  application_name,
  arn,
  compute_platform,
  linked_to_github,
  tags
from 
  aws_codedeploy_app
where 
  arn = '{{ output.resource_aka.value }}-dummy';