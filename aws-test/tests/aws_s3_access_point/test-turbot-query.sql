select
  akas,
  title
from
  aws.aws_s3_access_point
where
  name = '{{ resourceName }}';