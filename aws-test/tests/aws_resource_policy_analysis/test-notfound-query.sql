select
  r.name,
  r.role_id,
  pa.is_public,
  pa.allowed_principal_account_ids,
  pa.allowed_principals,
  pa.allowed_principal_services,
  pa.allowed_organization_ids,
  pa.public_access_levels,
  pa.shared_access_levels,
  pa.private_access_levels,
  pa.public_statement_ids,
  pa.shared_statement_ids,
  pa.access_level,
  r.arn
from
  aws_iam_role as r,
  aws_resource_policy_analysis as pa
where
  pa.account_id = r.account_id
  and pa.policy = r.assume_role_policy_std
  and r.name = '{{ resourceName }}-not-present';