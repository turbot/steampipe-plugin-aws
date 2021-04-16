select
  name,
  access_point_arn
from
  aws.aws_s3_access_point
where
  name = 'dummy-{{ resourceName }}'
  and region = '{{ output.region_id.value }}';