select
  namespace_name,
  namespace_arn,
  namespace_id,
  region
from
  aws_redshiftserverless_namespace
where
  namespace_name = '{{ output.resource_name.value }}-dummy';