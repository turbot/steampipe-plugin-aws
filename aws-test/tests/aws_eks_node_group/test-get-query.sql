select name, arn, version
from aws.aws_eks_cluster
where name = '{{output.node_group_name.value}}'
