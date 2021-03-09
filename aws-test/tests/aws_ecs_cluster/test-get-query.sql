select cluster_arn, cluster_name, active_sevices_count, status
from aws.aws_ecs_cluster
where cluster_arn = '{{ output.resource_id.value }}';