select fargate_profile_name, title, tags, akas, account_id
from aws.aws_eks_fargate_profile
where fargate_profile_name = 'dummy-{{ resourceName }}';