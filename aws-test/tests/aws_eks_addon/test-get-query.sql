select addon_name, addon_arn, addon_version, status
from aws.aws_eks_addon
where cluster_name = '{{ resourceName }}' and addon_name = 'vpc-cni';