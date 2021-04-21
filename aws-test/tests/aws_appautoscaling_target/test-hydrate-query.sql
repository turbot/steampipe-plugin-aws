select service_namespace, title
from aws.aws_appautoscaling_target
where service_namespace = '{{ output.resource_name.value }}';
