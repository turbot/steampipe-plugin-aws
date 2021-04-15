select
  name,
  access_point_arn,
  bucket_name,
  block_public_acls,
  block_public_policy,
  ignore_public_acls,
  restrict_public_buckets,
  network_origin,
  vpc_id,
  account_id,
  region,
  partition
from 
  aws.aws_s3_access_point
where 
  name = '{{ resourceName }}';
