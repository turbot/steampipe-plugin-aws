select tags
from aws_acmpca_certificate_authority
where arn = '{{ output.resource_aka.value }}'
