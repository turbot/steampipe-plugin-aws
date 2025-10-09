select cluster_name, principal_arn, policy_arn, access_scope_type, associated_at, modified_at
from aws.aws_eks_access_policy_association
where cluster_name = '{{ resourceName }}'
  and principal_arn = '{{ output.principal_arn.value }}'
  and policy_arn = '{{ output.policy_arn.value }}';

