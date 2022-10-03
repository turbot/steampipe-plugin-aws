select cluster_arn, cluster_name, akas, tags, title
from aws_ecs_cluster
where cluster_arn = '{{ output.resource_id.value }}';