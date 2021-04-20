select cluster_name, arn, tags, description, node_type
from aws.aws_dax_cluster
where cluster_name = '{{ output.resource_name.value }}';
