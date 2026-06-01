select
  registry_id,
  arn,
  replication_configuration -> 'Rules' -> 0 -> 'Destinations' -> 0 ->> 'Region' as destination_region,
  policy ->> 'Version' as policy_version,
  region
from aws.aws_ecr_registry
where region = '{{ output.aws_region.value }}';
