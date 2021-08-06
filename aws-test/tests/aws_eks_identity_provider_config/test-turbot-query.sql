select name, title, tags, akas, account_id
from aws.aws_eks_identity_provider_config
where arn = '{{ output.resource_aka.value }}';