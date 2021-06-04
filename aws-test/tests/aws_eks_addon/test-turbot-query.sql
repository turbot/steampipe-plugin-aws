select title, tags, akas, account_id
from aws.aws_eks_addon
where addon_arn = '{{ output.resource_aka.value }}';