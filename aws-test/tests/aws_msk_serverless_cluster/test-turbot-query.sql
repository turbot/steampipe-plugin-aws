select
  arn,
  akas,
  title,
  tags
from
  aws_msk_serverless_cluster
where
  cluster_name = '{{ resourceName }}';
