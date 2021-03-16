select name, arn, endpoint, status
from aws.aws_eks_cluster
where arn = '{{ output.resource_aka.value }}';
