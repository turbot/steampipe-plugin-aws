select account_id, region, title
from aws.aws_appautoscaling_target
where service_namespace = '{{ output.resource_name.value }}';
