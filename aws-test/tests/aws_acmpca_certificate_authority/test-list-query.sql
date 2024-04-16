select arn, title, akas
from aws_acmpca_certificate_authority
where akas::text = '["{{ output.resource_aka.value }}"]'
