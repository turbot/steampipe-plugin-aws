select addon_version, account_id
from aws.aws_eks_addon_version
where addon_version = '{{ output.addon_version.value }}' and region = '{{ output.region_id.value }}'  and addon_name = '{{ output.addon_name.value }}' ;