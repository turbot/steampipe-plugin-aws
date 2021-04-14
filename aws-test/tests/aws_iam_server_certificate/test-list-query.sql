select 
  name, 
  arn
from 
  aws.aws_iam_server_certificate
where 
  arn = '{{ output.resource_aka.value }}';
