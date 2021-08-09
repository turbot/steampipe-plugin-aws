select name, arn
from aws.aws_eks_identity_provider_config
where arn = '{{ output.resource_aka.value }}';
