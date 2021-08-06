select name, type, arn, cluster_name
from aws.aws_eks_identity_provider_config
where name = '{{ resourceName }}';
