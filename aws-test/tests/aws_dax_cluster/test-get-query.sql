select cluster_name, arn, tags
from aws.aws_dax_cluster
where cluster_name = '{{ output.resource_name.value }}';
