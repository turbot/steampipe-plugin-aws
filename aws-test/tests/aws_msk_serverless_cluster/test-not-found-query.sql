select
  cluster_name,
  arn,
  akas,
  serverless ->> 'VpcConfigs' as vpc_config,
  serverless ->> 'ClientAuthentication' as client_authentication
from
  aws_msk_serverless_cluster
where
  cluster_name = '{{ resourceName }}-dummy';
