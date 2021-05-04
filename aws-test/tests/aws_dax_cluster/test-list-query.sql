select akas, cluster_name, title
from aws.aws_dax_cluster
where cluster_name = '{{ output.resource_name.value }}';
