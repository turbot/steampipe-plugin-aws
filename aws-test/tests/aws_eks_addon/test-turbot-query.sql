select title, tags, akas, account_id
from aws.aws_eks_addon
where arn = '{{ output.resource_aka.value }}';