select certificate_arn, domain_name, in_use_by, issuer
from aws_acm_certificate
where certificate_arn = '{{ output.resource_aka.value }}'
