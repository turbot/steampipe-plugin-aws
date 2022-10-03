select
  name,
  access_point_policy_is_public,
  policy,
  policy_std
from
  aws.aws_s3_access_point
where
  name = '{{ resourceName }}'
  and region = '{{ output.region_id.value }}';
