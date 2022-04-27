select quota_name, quota_arn
from aws_servicequotas_default_service_quota
where quota_arn = '{{ output.resource_aka.value }}' and service_code = '{{ output.service_code.value }}';
