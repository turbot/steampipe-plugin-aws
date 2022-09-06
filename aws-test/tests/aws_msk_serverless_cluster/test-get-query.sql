select 
  cluster_name,
  cluster_arn,
  akas,
  tags
from 
  aws_msk_serverless_cluster
where 
  cluster_arn = '{{ output.resource_aka.value }}';