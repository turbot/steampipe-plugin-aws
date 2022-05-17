select
  arn,
  authorized_account_id,
  authorized_aws_region
from
  aab_aab.aws_config_aggregate_authorization
where akas::text = '["{{ output.resource_aka.value }}"]';