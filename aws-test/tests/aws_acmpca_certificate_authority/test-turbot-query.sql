select title, akas, region, account_id
from aws_acmpca_certificate_authority
where arn = '{{ output.resource_aka.value }}'
