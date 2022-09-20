select
  cluster_name,
  arn,
  akas,
  tags
from
  aws_msk_serverless_cluster
where
  cluster_name = '{{ resourceName }}';
