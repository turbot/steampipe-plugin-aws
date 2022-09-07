select 
  title,
  tags,
  akas,
  account_id
from 
  aws_codedeploy_app
where
  application_id = '{{ output.resource_id.value }}';
