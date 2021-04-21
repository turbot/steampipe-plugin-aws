select service_namespace, partition, region
from aws.aws_appautoscaling_target
where service_namespace = '{{ output.resource_name.value }}';
