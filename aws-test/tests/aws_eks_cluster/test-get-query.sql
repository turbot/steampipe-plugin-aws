select name, arn, endpoint, role_arn, version, resources_vpc_config, identity, status, certificate_authority -> 'Data'as certificate_authority_data
from aws.aws_eks_cluster
where arn = '{{ output.resource_aka.value }}'
