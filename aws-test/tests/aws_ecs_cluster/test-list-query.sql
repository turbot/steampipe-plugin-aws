select cluster_arn, cluster_name, active_services_count, attachments, attachments_status, capacity_providers, default_capacity_provider_strategy, pending_tasks_count,
registered_container_instances_count, running_tasks_count, status
from aws_ecs_cluster
where akas::text = '["{{ output.resource_aka.value }}"]';