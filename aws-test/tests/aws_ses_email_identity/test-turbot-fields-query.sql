select 
  title,
  akas, 
  region
from
  aws_ses_email_identity
where name = '{{resourceName}}';
