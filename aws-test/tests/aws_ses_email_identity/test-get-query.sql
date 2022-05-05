select
  name,
  arn,
  region,
  akas
from
  aws_ses_email_identity
where name = '{{resourceName}}';
