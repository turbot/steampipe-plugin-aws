select addon_name, arn, addon_version, status
from aws.aws_eks_addon
where cluster_name = '{{ resourceName }}' and addon_name = 'vpc-cni';