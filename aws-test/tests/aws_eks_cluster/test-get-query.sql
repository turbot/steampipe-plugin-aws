select name, arn, endpoint, role_arn, version, platform_version, status
from aws.aws_eks_cluster
where name = '{{ resourceName }}';
