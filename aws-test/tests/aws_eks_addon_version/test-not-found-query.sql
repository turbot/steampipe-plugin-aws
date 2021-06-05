select addon_version
from aws.aws_eks_addon_version
where addon_version = 'dummy-{{ output.addon_version.value }}';