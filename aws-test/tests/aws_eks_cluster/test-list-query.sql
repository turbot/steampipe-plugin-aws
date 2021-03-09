select name, arn, endpoint, status, role_arn
from aws.aws_eks_cluster
where arn = '{{ output.resource_aka.value }}'
