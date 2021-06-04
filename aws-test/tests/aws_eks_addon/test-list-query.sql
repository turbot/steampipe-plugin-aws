select addon_name, addon_arn, cluster_name, service_account_role_arn
from aws.aws_eks_addon
where addon_arn = '{{ output.resource_aka.value }}';