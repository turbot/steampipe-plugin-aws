select cluster_name, principal_arn, type
from aws.aws_eks_access_entry
where cluster_name = 'non-existent-cluster' and principal_arn = 'arn:aws:iam::123456789012:role/non-existent-role';

