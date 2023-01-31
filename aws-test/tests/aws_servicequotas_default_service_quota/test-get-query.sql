select quota_name, quota_arn
from aws_servicequotas_default_service_quota
where quota_code = '{{ output.quota_code.value }}' and service_code = '{{ output.service_code.value }}' and region = 'us-east-2';
