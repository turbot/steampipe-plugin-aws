select service_namespace, region, account_id
from aws.aws_appautoscaling_target
where service_namespace = '{{ output.resource_name.value }}3p';

