select
  namespace_id,
  title,
  tags,
  akas
from
  aws_redshiftserverless_namespace
where
  namespace_name = '{{ output.resource_name.value }}';