select
  name,
  cloud_watch_encryption,
  job_bookmarks_encryption,
  s3_encryption
from
  aws_glue_security_configuration
where 
  name = 'dummy-{{ resourceName }}';