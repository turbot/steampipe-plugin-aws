select title, akas, region, account_id
from aws_acm_certificate
where certificate_arn = '{{ output.resource_aka.value }}'
