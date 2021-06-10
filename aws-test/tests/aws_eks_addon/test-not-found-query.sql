select addon_name, arn, addon_version, status
from aws.aws_eks_addon
where cluster_name = 'dummy-{{ resourceName }}' and addon_name = 'dummy-vpc-cni';