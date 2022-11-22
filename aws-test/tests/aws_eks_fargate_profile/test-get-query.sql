select fargate_profile_name, cluster_name, fargate_profile_arn, status
from aws.aws_eks_fargate_profile
where fargate_profile_name = '{{ resourceName }}' and cluster_name = '{{ resourceName }}';
