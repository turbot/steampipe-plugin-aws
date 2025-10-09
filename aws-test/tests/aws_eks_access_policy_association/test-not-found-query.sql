select cluster_name, principal_arn, policy_arn
from aws.aws_eks_access_policy_association
where cluster_name = 'non-existent-cluster'
  and principal_arn = 'arn:aws:iam::123456789012:role/non-existent-role'
  and policy_arn = 'arn:aws:eks::aws:cluster-access-policy/NonExistentPolicy';

