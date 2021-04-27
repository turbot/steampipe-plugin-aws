select service_namespace, partition, region, resource_id
from aws.aws_appautoscaling_target
where resource_id = '{{ output.resource_id.value }}';;
