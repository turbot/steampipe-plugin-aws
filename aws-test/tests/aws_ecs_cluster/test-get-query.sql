select cluster_arn, cluster_name, active_services_count, status
from aws_ecs_cluster
where cluster_arn = '{{ output.resource_id.value }}';