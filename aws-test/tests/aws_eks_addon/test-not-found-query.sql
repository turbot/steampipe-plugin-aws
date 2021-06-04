select addon_name, addon_arn, addon_version, status
from aws.aws_eks_addon
where cluster_name = 'dummy-{{ resourceName }}' and addon_name = 'dummy-vpc-cni';