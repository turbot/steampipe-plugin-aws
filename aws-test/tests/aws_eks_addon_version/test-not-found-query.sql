select addon_version
from aws.aws_eks_addon_version
where addon_version = 'dummy-{{ output.addon_version.value }}' and region = '{{ output.region_id.value }}';