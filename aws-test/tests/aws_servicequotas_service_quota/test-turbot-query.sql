select title, akas
from aws_servicequotas_service_quota
where quota_arn = '{{ output.resource_aka.value }}' and service_code = '{{ output.service_code.value }}';
