select partition, region, resource_id
from aws.aws_appautoscaling_target
where service_namespace = 'dynamodb' and title = '{{ output.resource_id.value }}';
