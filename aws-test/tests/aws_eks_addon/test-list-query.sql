select addon_name, arn, cluster_name, service_account_role_arn
from aws.aws_eks_addon
where arn = '{{ output.resource_aka.value }}';