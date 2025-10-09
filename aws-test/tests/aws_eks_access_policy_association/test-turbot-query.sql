select
  cluster_name,
  principal_arn,
  policy_arn,
  title,
  akas
from aws.aws_eks_access_policy_association
where cluster_name = '{{ resourceName }}'
  and principal_arn = '{{ output.principal_arn.value }}'
  and policy_arn = '{{ output.policy_arn.value }}';

