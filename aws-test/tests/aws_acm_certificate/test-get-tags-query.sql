select tags
from aws_acm_certificate
where certificate_arn = '{{ output.resource_aka.value }}'
