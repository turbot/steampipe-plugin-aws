select cluster_name, akas, tags, title
from aws.aws_dax_cluster
where cluster_name = '{{ output.resource_name.value }}';
