select cluster_name, principal_arn, policy_arn, access_scope_type
from aws.aws_eks_access_policy_association
where cluster_name = '{{ resourceName }}' and principal_arn = '{{ output.principal_arn.value }}';

