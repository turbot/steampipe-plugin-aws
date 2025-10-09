select cluster_name, principal_arn, type, username
from aws.aws_eks_access_entry
where cluster_name = '{{ resourceName }}';

