select 
  name, 
  arn, 
  certificate_body,
  path,
  region,
  account_id
from 
  aws.aws_iam_server_certificate
where 
  name = '{{ resourceName }}';
