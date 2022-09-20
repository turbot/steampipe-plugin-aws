select
  cluster_name,
  arn,
  akas,
  tags
from
  aws_msk_serverless_cluster
where
  arn = '{{ output.resource_aka.value }}';
