select
  namespace_name,
  namespace_arn,
  namespace_id,
  region
from
  aws_redshiftserverless_namespace
where
  namespace_arn = '{{ output.resource_aka.value }}';