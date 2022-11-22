select fargate_profile_name, title, akas, account_id
from aws.aws_eks_fargate_profile
where fargate_profile_arn = '{{ output.resource_aka.value }}';